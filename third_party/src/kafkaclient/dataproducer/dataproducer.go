package dataproducer

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Shopify/sarama"

	"kafkaclient"
	"paas_lager/lager"
)

type Producer struct {
	logger   lager.Logger
	Client   sarama.Client
	Producer sarama.AsyncProducer
}

func NewProducerClient(conf *kafkaclient.KafkaConfig, logger lager.Logger) *Producer {

	clientConfig := sarama.NewConfig()

	clientConfig.ClientID = kafkaclient.GenerateGuid()
	clientConfig.Producer.RequiredAcks = sarama.NoResponse
	clientConfig.Producer.Timeout = 1000 * time.Millisecond
	clientConfig.Producer.Return.Successes = true

	client, err := sarama.NewClient(conf.KafkaBrokers, clientConfig)
	if err != nil {
		logger.Error("[dataproducer] connect error!", err, lager.Data{"kafkaBrokers": conf.KafkaBrokers})
		return nil
	}

	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		logger.Error("[dataproducer] create produce error!", err)
		return nil
	}

	return &Producer{
		Client:   client,
		Producer: producer,
		logger:   logger,
	}
}

func (producerClient *Producer) ProduceMsg(topic string, message interface{}) error {

	msg, err := json.Marshal(message)
	if err != nil {
		producerClient.logger.Error("[dataproducer] Marshal error!", err)
		return err
	}

//当kafka机器异常的时候，写消息到input channel会出现阻塞，不会继续往下走，所以在这里加了超时处理
produceMsgLoop:
	for {
		select {
			case producerClient.Producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.ByteEncoder(msg)}:
				break produceMsgLoop
			case <-time.After(1000 * time.Millisecond):
				break produceMsgLoop
		}

	}
	producerClient.logger.Debug("Produce message success!")
	return producerClient.CheckErrors()
}

func (producerClient *Producer) ProduceMsgByteArray(topic string, data []byte) error {
//当kafka机器异常的时候，写消息到input channel会出现阻塞，不会继续往下走，所以在这里加了超时处理
produceMsgByteArrayLoop:
        for {
                select {
			case producerClient.Producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.ByteEncoder(data)}:
				break produceMsgByteArrayLoop
			case <-time.After(1000 * time.Millisecond):
				break produceMsgByteArrayLoop
		}
	}
	producerClient.logger.Debug("Produce message byte array  success")
	return producerClient.CheckErrors()
}

func (producerClient *Producer) CheckErrors() error {

	var err error
	logger := producerClient.logger
	for {
		select {
		case err := <-producerClient.Producer.Errors():
			logger.Error("[dataproducer] producer error!", err)
			return err
		case <-producerClient.Producer.Successes():
			logger.Debug("[dataproducer] producer Success!")
			return nil
		case <-time.After(1000 * time.Millisecond):
			logger.Error("[dataproducer] producer error!", errors.New("producer send message timeout!"))
			return err
		}
	}
}

func (producerClient *Producer) Close() error {
	err := producerClient.Producer.Close()
	if err != nil {
		return err
	}

	err = producerClient.Client.Close()
	if err != nil {
		return err
	}
	return nil
}
