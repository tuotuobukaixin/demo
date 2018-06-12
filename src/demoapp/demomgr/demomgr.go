package main

import (
	"math/rand"
	"net/http"
	"os"
	"demoapp/demomgr/service"
	"demoapp/demomgr/conf"
	"demoapp/common/models"
	"demoapp/common"
	"time"
)

// @APIVersion 2.0.0
// @APITitle kubernetes-runtime API
// @APIDescription Our API usually works as expected.
//

const (
	DefaultReadTimeout  = 5 * time.Second
	DefaultWriteTimeout = 60 * time.Second
)

func main() {
	common.LOGGER.Info("Cron task start...")


	router, err := service.GetRouter()
	if err != nil {
		common.LOGGER.Error("GetRouter", err)
		os.Exit(1)
	}
	http.Handle("/", router)
	models.Setup(map[string]string{"DatasourceURL": conf.Config.DatasourceURL})


	//err = http.ListenAndServe(hostIP+":"+port, nil)
	server := &http.Server{Addr: "0.0.0.0:8088", ReadTimeout: DefaultReadTimeout, WriteTimeout: DefaultWriteTimeout}
	err = server.ListenAndServe()
	if err != nil {
		common.LOGGER.Error("ListenAndServe", err)
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	os.Exit(0)
}
