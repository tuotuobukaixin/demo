package main

import (
	"eventclient"
	"eventclient/models"
	"flag"
	"fmt"

	"kafkaclient"
	"kafkaclient/dataproducer"
	"paas_lager"
	"time"
)

func main() {
	flag.Parse()

	paas_lager.Init(paas_lager.Config{
		LoggerLevel:   "ERROR",
		LoggerFile:    "./eventtest.log",
		LogFormatText: false,
		EnableRsyslog: true,
	})
	logger := paas_lager.NewLogger("eventbus_connect")

	event := models.EventMessage{
		ObjectId:          "AppID=App124235,NodeID=node-test1,InstanceID=ins1234",
		ObjectClass:       "Application",
		SerialNumber:      10000000000001,
		EventTime:         "2015-04-15T14:12:56Z",
		EventType:         12,
		EventId:           100001,
		EventName:         "VMAppInsCreate",
		PerceivedSeverity: "Info",
	}

	alarm := models.AlarmMessage{
		ObjectId:          "AppID=App124235,NodeID=node-test1,InstanceID=ins1234",
		ObjectClass:       "Application",
		SerialNumber:      10000000000001,
		EventTime:         "2015-04-15T14:12:56Z",
		ClearedType:       "ADMC",
		EventType:         7,
		EventId:           100005,
		EventName:         "VMAppInsStartFail",
		PerceivedSeverity: "Major",
	}

	//创建producer
	var conf kafkaclient.KafkaConfig
	conf.KafkaBrokers = append(conf.KafkaBrokers, "localhost:9092")
	producer := dataproducer.NewProducerClient(&conf, logger)

	//最后的参数1代表部署子系统
	clientGroup := eventclient.NewDefaultClient(logger, producer, 1)

	if clientGroup == nil {
		logger.Info("eventclient create failed!")
	}

	for {
		fmt.Println("Publish event message", event)
		clientGroup.EventPublish(event)
		fmt.Println("Publish alarm message", alarm)
		clientGroup.AlarmPublish(alarm)
		time.Sleep(10 * time.Second)
	}
}
