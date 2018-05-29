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
	var tmp = models.Job{}
	tmp.Name = util.Config.Jobname
	tmp.Status = "Running"
	models.AddJob(&tmp)
	time.Sleep(time.Duration(util.Config.Timeout) * time.Second)
	tmp.Status= "Success"
	models.UpdateJob(&tmp)
	os.Exit(0)
}
