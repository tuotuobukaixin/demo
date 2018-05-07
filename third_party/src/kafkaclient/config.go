package kafkaclient

import (
	uuid "github.com/nu7hatch/gouuid"
)

type KafkaConfig struct {
	Groupid        string         `json:"groupid"`
	KafkaBrokers   []string       `json: kafkaBrokers`
	ZookeeprConfig ZookeeprConfig `json:"zookeeperConfig"`
}

type ZookeeprConfig struct {
	Addresses         []string `json:"addresses"`
	RootDir           string   `json:"rootDir"`
	Timeout           int      `json:"timeout"` //second
	MaxRequestRetries int      `json:"maxRequestRetries"`
	RequestBackoff    int      `json:"requestBackoff"` //second
}

func GenerateGuid() string {
	uuid, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return uuid.String()
}
