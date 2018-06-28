package util

import (
	"encoding/json"
	"io/ioutil"
)

type configuration struct {
	Num   int
	Redisserver string
	Timeout int
}

// The configuration structure exposed to the others to get the configuration information
// Tht gcrypto engine
var Config configuration

func init() {
	content, err := ioutil.ReadFile("app.conf")
	if err == nil {
		err = json.Unmarshal(content, &Config)
		if err != nil {
			panic(err)
		}
	}

}
