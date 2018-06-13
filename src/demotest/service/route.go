package service

import (
	"net/http"

	"github.com/tedsuo/rata"
	"demotest/util"
)


var rteHandlers = rata.Handlers{
	"get_demotest":    http.HandlerFunc(GetGameserverDetail),
	"health":    http.HandlerFunc(Health),
}

func GetRouter() (router http.Handler, err error) {
	var rteRoutes = rata.Routes{
		{Name: "get_demotest", Method: "GET", Path: "/api/v1/"+util.Config.ServerName},
		{Name: "health", Method: "GET", Path: "/api/v1/demotesthealth"},
	}
	router, err = rata.NewRouter(rteRoutes, rteHandlers)
	return
}
