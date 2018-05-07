package dataproducer_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"kafkaclient"
	"kafkaclient/datacollector"
	. "kafkaclient/dataproducer"
	"paas_lager"
)

var _ = Describe("Dataproducer", func() {

	paas_lager.Init(paas_lager.Config{
		LoggerLevel:   "DEBUG",
		LoggerFile:    "./dataproducer.log",
		EnableRsyslog: true,
	})
	logger := paas_lager.NewLogger("producer")

	var (
		localBroker = "localhost:9092"
		localZk     = "localhost:2181"
		topic       string
		conf        kafkaclient.KafkaConfig
		collector   *datacollector.Collector
	)

	BeforeEach(func() {
		groupid := fmt.Sprintf("test-groupid-%d", time.Now().Unix())
		topic = fmt.Sprintf("test-topic-%d", time.Now().Unix())
		addresses := []string{localZk}
		zookeeprConfig := kafkaclient.ZookeeprConfig{
			Addresses:         addresses,
			RootDir:           "",
			Timeout:           10,
			MaxRequestRetries: 3,
			RequestBackoff:    150,
		}

		kafkaBrokers := []string{localBroker}
		conf = kafkaclient.KafkaConfig{
			Groupid:        groupid,
			KafkaBrokers:   kafkaBrokers,
			ZookeeprConfig: zookeeprConfig,
		}

		collector = datacollector.NewCollector(&conf)
		Expect(collector).NotTo(BeNil())

	})

	Describe("Create", func() {
		Context("Create KAFKA producer", func() {
			It("Create KAFKA producer Success", func() {
				producer := NewProducerClient(&conf, logger)
				Expect(producer).NotTo(BeNil())
			})

		})

	})
	Describe("Send Msg", func() {
		Context("Send Msg", func() {
			It("Send Msg Success", func() {
				producer := NewProducerClient(&conf, logger)
				Expect(producer).NotTo(BeNil())

				By("create topic")
				produceN(1, topic, localBroker, "")
				fmt.Println(topic)
				time.Sleep(2 * time.Second)

				By("Start Collector Msg!")
				collector.Run([]string{topic}, 1)
				msgflag := false
				go func() {
					for {
						select {
						case msg := <-collector.MessageChan:
							By("Get Message !!!")
							fmt.Println("msgNonAlarm", string(msg))
							Expect(string(msg)).NotTo(Equal(""))
							msgflag = true
							break
						}
					}
				}()

				time.Sleep(2 * time.Second)
				By("Producer Send Message!!!")

				for i := 0; i < 3; i++ {
					By(" Send Message!")
					data := "metrics.DataPoint{}"
					err := producer.ProduceMsg(topic, data)
					Expect(err).To(BeNil())
					time.Sleep(1 * time.Second)
				}

				By("Check Msg!")
				Expect(msgflag).To(Equal(true))

				err := producer.Close()
				Expect(err).To(BeNil())

			})

		})

	})
})

func produceN(n int, topic string, brokerAddr string, message interface{}) {
	clientConfig := sarama.NewConfig()
	clientConfig.Producer.Timeout = 10 * time.Second
	client, err := sarama.NewClient([]string{brokerAddr}, clientConfig)
	Expect(err).ShouldNot(HaveOccurred())

	defer client.Close()
	msg, err := json.Marshal(message)
	Expect(err).ShouldNot(HaveOccurred())
	producer, err := sarama.NewAsyncProducerFromClient(client)
	Expect(err).ShouldNot(HaveOccurred())

	defer producer.Close()
	for i := 0; i < n; i++ {
		producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.ByteEncoder(msg)}
	}
	select {
	case e := <-producer.Errors():
		Expect(e).ShouldNot(HaveOccurred())
	case <-time.After(5 * time.Second):
	}
}
