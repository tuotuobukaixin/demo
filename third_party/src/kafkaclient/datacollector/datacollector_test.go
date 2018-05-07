package datacollector_test

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"policy_engine/config"
	. "policy_engine/datamanager/datacollector"

	"github.com/Shopify/sarama"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	kafkaClient "github.com/stealthly/go_kafka_client"
)

var localZk = "localhost:2181"
var localBroker = "localhost:9092"

var _ = Describe("Datacollector", func() {

	Describe("Kafka Connect Success", func() {
		It("Connect Success", func() {
			//rootDir := fmt.Sprintf("test-dir-%d", time.Now().Unix())
			groupid := fmt.Sprintf("test-groupid-%d", time.Now().Unix())
			topic := fmt.Sprintf("test-topic-%d", time.Now().Unix())
			produceMessages := 3

			addresses := []string{localZk}
			zookeeprConfig := config.ZookeeprConfig{
				Addresses:         addresses,
				RootDir:           "",
				Timeout:           10,
				MaxRequestRetries: 3,
				RequestBackoff:    150,
			}

			conf := config.KafkaConfig{
				Groupid:        groupid,
				ZookeeprConfig: zookeeprConfig,
			}

			collector := NewCollector(&conf)
			Expect(collector).NotTo(BeNil())

			fmt.Println(topic)
			produceN(1, topic, localBroker)
			//CreateMultiplePartitionsTopic(localZk, topic, 1)
			//EnsureHasLeader(&zookeeprConfig, topic)
			time.Sleep(1 * time.Second)
			go collector.Consumer.StartStatic(map[string]int{topic: 1})
			time.Sleep(1 * time.Second)
			By("Producer Send Message!!!")
			go produceN(produceMessages, topic, localBroker)
			By("Waiting for some seconds before producing another message")

			msgflag := false

			go func() {
				for {
					select {

					case actual := <-collector.MessageChan:
						By("Get Message !!!")
						fmt.Println("msgNonAlarm", string(actual))
						Expect(string(actual)).NotTo(Equal(""))
						msgflag = true
						break
					}
				}
			}()

			time.Sleep(5 * time.Second)

			closeWithin(5*time.Second, collector.Consumer)
			//closeWithin(20*time.Second, collector.Consumer)
			Expect(msgflag).To(Equal(true))
			By("Kafka Start Success!!!")

		})
	})
})

func produceN(n int, topic string, brokerAddr string) {
	clientConfig := sarama.NewConfig()
	clientConfig.Producer.Timeout = 10 * time.Second
	client, err := sarama.NewClient([]string{brokerAddr}, clientConfig)
	Expect(err).ShouldNot(HaveOccurred())

	defer client.Close()

	producer, err := sarama.NewAsyncProducerFromClient(client)
	Expect(err).ShouldNot(HaveOccurred())

	defer producer.Close()
	for i := 0; i < n; i++ {
		producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder(fmt.Sprintf("test-kafka-message-%d", i))}
	}
	select {
	case e := <-producer.Errors():
		Expect(e).ShouldNot(HaveOccurred())
	case <-time.After(5 * time.Second):
	}
}

func closeWithin(timeout time.Duration, consumer *kafkaClient.Consumer) {
	select {
	case <-consumer.Close():
		{
			fmt.Println("Success!!")
		}
	case <-time.After(timeout):
		{
			fmt.Println("Fail!!")
		}
	}
}

//Convenience utility to create a topic topicName with numPartitions partitions in Zookeeper located at zk (format should be host:port).
//Please note that this requires Apache Kafka 0.8.1 binary distribution available through KAFKA_PATH environment variable
func CreateMultiplePartitionsTopic(zk string, topicName string, numPartitions int) {
	if runtime.GOOS == "windows" {
		params := fmt.Sprintf("--create --zookeeper %s --replication-factor 1 --partitions %d --topic %s", zk, numPartitions, topicName)
		script := fmt.Sprintf("%s\\bin\\windows\\kafka-topics.bat %s", os.Getenv("KAFKA_PATH"), params)
		exec.Command("cmd", "/C", script).Output()
	} else {
		params := fmt.Sprintf("--create --zookeeper %s --replication-factor 1 --partitions %d --topic %s", zk, numPartitions, topicName)
		script := fmt.Sprintf("%s/bin/kafka-topics.sh %s", os.Getenv("KAFKA_PATH"), params)
		out, err := exec.Command("sh", "-c", script).Output()
		if err != nil {
			panic(err)
		}
		fmt.Println("create topic", out)
	}
}

//blocks until the leader for every partition of a given topic appears
//this is used by tests only to avoid "In the middle of a leadership election, there is currently no leader for this partition and hence it is unavailable for writes"
/*func EnsureHasLeader(conf *config.ZookeeprConfig, topic string) {

	zookeeper := kafkaClient.NewZookeeperCoordinator(CreateZookeeperConf(conf))

	zookeeper.Connect()
	hasLeader := false

	numPartitions := 0
	for !hasLeader {
		var topicInfo *kafkaClient.TopicInfo
		var err error
		for i := 0; i < 3; i++ {
			topicInfo, err = zookeeper.getTopicInfo(topic)
			if topicInfo != nil {
				break
			}
		}
		if err != nil {
			continue
		}
		numPartitions = len(topicInfo.Partitions)

		hasLeader = true
		for partition, leaders := range topicInfo.Partitions {
			if len(leaders) == 0 {
				fmt.Println("Partition has no leader, waiting...", partition, zookeeper)
				hasLeader = false
				break
			}
		}

		if !hasLeader {
			time.Sleep(1 * time.Second)
		}
	}
	time.Sleep(time.Duration(numPartitions) * time.Second)
}*/
