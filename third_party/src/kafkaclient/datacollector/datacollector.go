package datacollector

import (
	"fmt"
	"sync"
	"time"

	kafkaClient "github.com/stealthly/go_kafka_client"
	"kafkaclient"
)

var consumedMessagesLock sync.Mutex

type Collector struct {
	MessagesLock *sync.Mutex
	Consumer     *kafkaClient.Consumer
	Conf         *kafkaClient.ConsumerConfig
	MessageChan  chan []byte
}

func NewCollector(conf *kafkaclient.KafkaConfig) *Collector {

	consumerConf := kafkaClient.DefaultConsumerConfig()
	consumerConf.Groupid = conf.Groupid
	consumerConf.Clientid = kafkaclient.GenerateGuid()
	consumerConf.Coordinator = kafkaClient.NewZookeeperCoordinator(CreateZookeeperConf(&conf.ZookeeprConfig))

	// consumer获取消息失败情况下，执行的回调函数。
	consumerConf.WorkerFailureCallback = func(_ *kafkaClient.WorkerManager) kafkaClient.FailedDecision {
		return kafkaClient.CommitOffsetAndContinue
	}
	consumerConf.WorkerFailedAttemptCallback = func(_ *kafkaClient.Task, _ kafkaClient.WorkerResult) kafkaClient.FailedDecision {
		return kafkaClient.CommitOffsetAndContinue
	}

	messageChan := make(chan []byte)
	consumerConf.Strategy = CreateKafkaStrategy(messageChan)
	consumer := kafkaClient.NewConsumer(consumerConf)

	return &Collector{
		Consumer:     consumer,
		MessagesLock: &consumedMessagesLock,
		Conf:         consumerConf,
		MessageChan:  messageChan,
	}
}

func (c *Collector) Run(topics []string, goroutinesNum int) {
	topicMap := make(map[string]int)
	for _, topic := range topics {
		topicMap[topic] = goroutinesNum
	}

	go c.Consumer.StartStatic(topicMap)
}

func (c *Collector) Close() {
	select {
	case <-c.Consumer.Close():
		{
			fmt.Println("Close kafka Collector Success!")
		}
	case <-time.After(2 * time.Second):
		{
			fmt.Println("Close kafka Collector Fail!")
		}
	}
}

func inLock(lock *sync.Mutex, fun func()) {
	lock.Lock()
	defer lock.Unlock()

	fun()
}

func CreateKafkaStrategy(messageChan chan []byte) kafkaClient.WorkerStrategy {
	return func(_ *kafkaClient.Worker, msg *kafkaClient.Message, id kafkaClient.TaskId) kafkaClient.WorkerResult {
		message := []byte(msg.Value)
		inLock(&consumedMessagesLock, func() {
			messageChan <- message
		})
		return kafkaClient.NewSuccessfulResult(id)
	}
}

func CreateZookeeperConf(conf *kafkaclient.ZookeeprConfig) *kafkaClient.ZookeeperConfig {
	return &kafkaClient.ZookeeperConfig{
		ZookeeperConnect:  conf.Addresses,
		ZookeeperTimeout:  time.Duration(conf.Timeout) * time.Second,
		MaxRequestRetries: conf.MaxRequestRetries,
		RequestBackoff:    time.Duration(conf.RequestBackoff) * time.Second,
		Root:              conf.RootDir,
	}
}
