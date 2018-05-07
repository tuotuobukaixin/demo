eventbus client SDK
===========

Eventbus Client SDK

---
## 使用说明

### import引用说明
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
eventclient是主要的引用包
eventclient/models提供消息定义
paas_lager日志输出模块
kafkaclient连接kafka的配置
kafkaclient/dataproducer产生kafka的producer

### 依赖第三方库
gopkg.in/validator.v2，参数有效性检测库
需要在使用该SDK时，使用go get 命令获取以上的第三方库或者直接引用版本内的库

### client使用说明

封装上报消息，提供了两种消息格式，使用方式如下：

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
        EventType:         16,
        EventId:           100005,
        EventName:         "VMAppInsCreateFail",
        PerceivedSeverity: "Major",
    }

创建producer对象实例
    //创建producer
    var conf kafkaclient.KafkaConfig
    conf.KafkaBrokers = append(conf.KafkaBrokers, "localhost:9092")
    producer := dataproducer.NewProducerClient(&conf, logger)

创建client对象实例

    //最后的参数1代表部署子系统
    clientGroup := eventclient.NewDefaultClient(logger, producer, 1)

    if clientGroup == nil {
        logger.Info("eventclient create failed!")
    }

发送事件消息

    fmt.Println("Publish event message", event)
    clientGroup.EventPublish(event)
    fmt.Println("Publish alarm message", alarm)
    clientGroup.AlarmPublish(alarm)
    time.Sleep(10 * time.Second)


启用开发者模式
    取消  "kafkaclient"，"kafkaclient/dataproducer"的引入
    将日志级别改成DEBUG
    取消创建producer的操作
    创建client对象实例的时候producer请传入nil

demo 代码例子请参考eventclient/eventstest


