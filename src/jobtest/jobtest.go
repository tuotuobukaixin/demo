package main

import (
	"os"
	"jobtest/util"
	"jobtest/models"
	"time"
)

// @APIVersion 2.0.0
// @APITitle kubernetes-runtime API
// @APIDescription Our API usually works as expected.
//

func theard() {

	for  {
		_=0
	}

}


func main() {

	models.Setup(map[string]string{"DatasourceURL": util.Config.DatasourceURL})

	go theard()

	time.Sleep(time.Duration(util.Config.Timeout) * time.Second)
	models.AddJob(util.Config.Jobname)
	os.Exit(0)
}
