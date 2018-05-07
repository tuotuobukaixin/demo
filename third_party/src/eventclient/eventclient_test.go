package eventclient_test

import (
	. "eventclient"
	"eventclient/models"
	"flag"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"kafkaclient"
	"kafkaclient/dataproducer"
	"paas_lager"
)

var _ = Describe("Eventclient", func() {
	flag.Parse()

	paas_lager.Init(paas_lager.Config{
		LoggerLevel:   "ERROR",
		LoggerFile:    "./eventtest.log",
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
	Describe("Loadsystemid", func() {
		Context("Loadsystemid", func() {
			It("should be success to Loadsystemid", func() {

				Expect(producer).NotTo(BeNil())
				clientGroup := NewDefaultClient(logger, producer, 1)
				Expect(clientGroup).NotTo(BeNil())
				fmt.Println("Publish event message", event)
				err := clientGroup.EventPublish(event)
				Expect(err).NotTo(HaveOccurred())
				fmt.Println("Publish alarm message", alarm)
				err = clientGroup.AlarmPublish(alarm)
				Expect(err).NotTo(HaveOccurred())

			})
		})
	})

})
