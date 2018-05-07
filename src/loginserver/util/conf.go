package util

import (
	"encoding/json"
	"io/ioutil"
)

type configuration struct {
	Registerurl   string
	ServerAddr    string
	DatasourceURL string
	Httpport      string
}

// The configuration structure exposed to the others to get the configuration information
// Tht gcrypto engine
var Config configuration

func init() {
	content, err := ioutil.ReadFile("conf/app.conf")
	if err == nil {
		err = json.Unmarshal(content, &Config)
		if err != nil {
			panic(err)
		}
	}

	Config.DatasourceURL = "root:root@tcp(10.186.69.39:3306)/game?charset=utf8&loc=Asia%2FShanghai"
	Config.Registerurl = "http://127.0.0.1:8087"
	Config.ServerAddr = "10.25.125.11"
	Config.Httpport = "8089"

}
