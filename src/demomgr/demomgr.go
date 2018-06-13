package main

import (
	"errors"
	"math/rand"
	"net/http"
	"os"
	"demomgr/service"
	"demomgr/util"
	"demomgr/models"
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
	util.LOGGER.Info("Cron task start...")


	router, err := service.GetRouter()
	if err != nil {
		util.LOGGER.Error("GetRouter", err)
		os.Exit(1)
	}
	http.Handle("/", router)
	models.Setup(map[string]string{"DatasourceURL": util.Config.DatasourceURL})

	hostIP := "0.0.0.0"
	if hostIP == "" {
		util.LOGGER.Error("HostIP null", errors.New("HostIP null"))
		os.Exit(1)
	}

	port := util.Config.Httpport
	if port == "" {
		port = "8087"
		util.LOGGER.Info("Listening in default port: " + port)
	} else {
		util.LOGGER.Info("Listening in port:" + port)
	}

	//err = http.ListenAndServe(hostIP+":"+port, nil)
	server := &http.Server{Addr: hostIP + ":" + port, ReadTimeout: DefaultReadTimeout, WriteTimeout: DefaultWriteTimeout}
	err = server.ListenAndServe()
	if err != nil {
		util.LOGGER.Error("ListenAndServe", err)
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	os.Exit(0)
}
