package sinks

import (
	"gopkg.in/validator.v2"
	"time"

	"eventclient/models"
	"kafkaclient/dataproducer"
	"paas_lager/lager"
)

const (
	ALARM_TOPIC = "alarm_message"
	EVENT_TOPIC = "event_message"
)

//A type of client:publish by kafka
type KafkaClient struct {
	logger        lager.Logger
	AlarmProducer *dataproducer.Producer
	EventProducer *dataproducer.Producer
	SysConfig     *models.SysConfig
}

func NewKafkaClient(
	logger lager.Logger,
	producer *dataproducer.Producer,
	sysConfig *models.SysConfig,
) *KafkaClient {

	eventClient := KafkaClient{
		logger:        logger,
		AlarmProducer: producer,
		EventProducer: producer,
		SysConfig:     sysConfig,
	}

	return &eventClient
}

func (client *KafkaClient) Close() {
	client.AlarmProducer.Close()
	client.EventProducer.Close()
}

func (client *KafkaClient) EventPublish(event models.EventMessage) error {
	logs := client.logger

	client.SysConfig.EventSerialNumber++
	t := time.Now().UTC()
	event.EventTime = t.Format(time.RFC3339)
	event.SerialNumber = client.SysConfig.EventSerialNumber

	err := validator.Validate(event)
	if err != nil {
		logs.Error("EventPublish event json.Marshal Fail!", err, lager.Data{"event": event})
		return err
	}

	return client.EventProducer.ProduceMsg(EVENT_TOPIC, event)
}

func (client *KafkaClient) AlarmPublish(alarm models.AlarmMessage) error {
	logs := client.logger

	client.SysConfig.AlarmSerialNumber++
	t := time.Now().UTC()
	alarm.EventTime = t.Format(time.RFC3339)
	alarm.SerialNumber = client.SysConfig.AlarmSerialNumber

	err := validator.Validate(alarm)
	if err != nil {
		logs.Error("AlarmPublish alarm json.Marshal Fail!", err, lager.Data{"alarm": alarm})
		return err
	}

	return client.AlarmProducer.ProduceMsg(ALARM_TOPIC, alarm)
}
