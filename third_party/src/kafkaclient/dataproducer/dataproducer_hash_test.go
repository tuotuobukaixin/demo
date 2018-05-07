package dataproducer_test

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"kafkaclient"
	"kafkaclient/datacollector"
	. "kafkaclient/dataproducer"
	"paas_lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DataproducerHash", func() {
	paas_lager.Init(paas_lager.Config{
		LoggerLevel:   "DEBUG",
		LoggerFile:    "./dataproducer.log",
		EnableRsyslog: true,
	})
	logger := paas_lager.NewLogger("producer")

	var (
		localBroker  = "localhost:9092"
		localZk      = "localhost:2181"
		topic        string
		conf         kafkaclient.KafkaConfig
		collector    *datacollector.Collector
		collector2   *datacollector.Collector
		hashproducer *HashProducer
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
		collector2 = datacollector.NewCollector(&conf)
		Expect(collector2).NotTo(BeNil())
		hashproducer = NewHashProducer(&conf, logger)
		Expect(hashproducer).NotTo(BeNil())

	})

	Describe("Send Hash Msg", func() {
		Context("Send Hash Msg", func() {
			It("Send Hash Msg Success", func() {
				type DataPoint struct {
					Metric    string            `json:"metric"`
					Timestamp int64             `json:"timestamp"`
					Value     interface{}       `json:"value"`
					Tags      map[string]string `json:"tags"`
				}
				tags := make(map[string]string)
				tags["app_guid"] = "test-app"
				tags["node_id"] = "test-node"

				tags2 := make(map[string]string)
				tags2["app_guid"] = "test-app"
				tags2["node_id"] = "test-node-22"
				data := DataPoint{
					Metric:    "cpu_usage",
					Timestamp: 121441,
					Value:     123,
					Tags:      tags,
				}
				By("create topic")
				hashproducer.SendMessage(topic, "app_guid:test-app_node_id:test-node", data)
				fmt.Println(topic)
				time.Sleep(2 * time.Second)

				By("Start Collector Msg!")
				collector.Run([]string{topic}, 1)
				collector2.Run([]string{topic}, 1)

				msgflag := false

				//判断同一个应用的消息不能出现在两个Consumer中
				getResultMsg := make(map[string]string)
				go func() {
					for {
						select {
						case msg := <-collector.MessageChan:
							By("Collector Get Message !!!")
							fmt.Println("msg", string(msg))
							Expect(string(msg)).NotTo(Equal(""))
							msgflag = true

							By("Judge Hash Message")
							var dataPoint DataPoint
							err := json.Unmarshal(msg, &dataPoint)
							Expect(err).ShouldNot(HaveOccurred())

							collectName, exist := getResultMsg[GenerateTagsStr(dataPoint.Tags)]
							if exist {
								Expect(collectName).To(Equal("Collector"))
							} else {
								getResultMsg[GenerateTagsStr(dataPoint.Tags)] = "Collector"
							}
							break
						case msg := <-collector2.MessageChan:
							By("Collector2 Get Message !!!")
							fmt.Println("msg", string(msg))
							Expect(string(msg)).NotTo(Equal(""))
							msgflag = true

							By("Judge Hash Message")
							var dataPoint DataPoint
							err := json.Unmarshal(msg, &dataPoint)
							Expect(err).ShouldNot(HaveOccurred())

							collectName, exist := getResultMsg[GenerateTagsStr(dataPoint.Tags)]
							if exist {
								Expect(collectName).To(Equal("Collector2"))
							} else {
								getResultMsg[GenerateTagsStr(dataPoint.Tags)] = "Collector2"
							}
							break
						}
					}
				}()

				time.Sleep(2 * time.Second)

				By("Producer Send Message!!!")
				key1 := GenerateTagsStr(tags)
				key2 := GenerateTagsStr(tags2)
				for i := 0; i < 4; i++ {
					By("send hash message")
					data.Tags = tags
					data.Timestamp = time.Now().Unix()
					hashproducer.SendMessage(topic, key1, data)
					time.Sleep(1 * time.Second)

					data.Tags = tags2
					data.Timestamp = time.Now().Unix()
					hashproducer.SendMessage(topic, key2, data)
					time.Sleep(1 * time.Second)
				}

				time.Sleep(3 * time.Second)

				By("Check Msg!")
				Expect(msgflag).To(Equal(true))

				hashproducer.Close()

			})

		})

	})
})

func GenerateTagsStr(tags map[string]string) string {
	sorted_keys := make([]string, 0)
	for key, _ := range tags {
		sorted_keys = append(sorted_keys, key)
	}
	sort.Strings(sorted_keys)
	tagStr := fmt.Sprintf("%s:%s", sorted_keys[0], tags[sorted_keys[0]])

	for i := 1; i < len(sorted_keys); i++ {
		if tags[sorted_keys[i]] != "" {
			singleTag := fmt.Sprintf("%s:%s", sorted_keys[i], tags[sorted_keys[i]])
			tagStr = fmt.Sprintf("%s-%s", tagStr, singleTag)
		}
	}

	return tagStr
}
