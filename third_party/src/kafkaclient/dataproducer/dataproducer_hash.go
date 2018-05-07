package dataproducer

import (
	"encoding/json"
	"errors"
	"kafkaclient"
	"paas_lager/lager"
	"time"

	kc "github.com/stealthly/go_kafka_client"
)

type HashProducer struct {
	logger   lager.Logger
	producer kc.Producer
}

func NewHashProducer(conf *kafkaclient.KafkaConfig, logger lager.Logger) *HashProducer {
	producerConfig := kc.DefaultProducerConfig()
	producerConfig.BrokerList = conf.KafkaBrokers
	producerConfig.Clientid = kafkaclient.GenerateGuid()

	//use hash partitioner
	producerConfig.Partitioner = kc.NewHashPartitioner

	producer := kc.NewSaramaProducer(producerConfig)
	return &HashProducer{
		logger:   logger,
		producer: producer,
	}

}

//key will be used to hash,get hash id.message will send to spec partition by the hash id.
//if you want to make sure message send success,you can call the HashProducer.producer.Error()
func (h *HashProducer) SendMessage(topic, key string, message interface{}) {
	producerMsg, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("msg-marshal-fail!", err, lager.Data{"key": key, "msg": message})
		return
	}

	producerMessage := &kc.ProducerMessage{
		Topic: topic,
		Key:   []byte(key),
		Value: producerMsg,
	}
//当kafka机器异常的时候，写消息到input channel会出现阻塞，不会继续往下走，所以在这里加了超时处理
	select {
	case h.producer.Input() <- producerMessage:
	case <-time.After(1000 * time.Millisecond):
		h.logger.Error("[dataproducer] producer error!", errors.New("producer send message timeout!"))
	}
	return
}

func (h *HashProducer) Close() {
	h.producer.Close()
}
