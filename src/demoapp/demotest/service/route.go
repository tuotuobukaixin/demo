package service

import (
	"net/http"

	"github.com/tedsuo/rata"
	"gameserver/util"
)


var rteHandlers = rata.Handlers{
	"get_gameserver":    http.HandlerFunc(GetGameserverDetail),
	"health":    http.HandlerFunc(Health),
}

func GetRouter() (router http.Handler, err error) {
	var rteRoutes = rata.Routes{
		{Name: "get_gameserver", Method: "GET", Path: "/api/v1/"+util.Config.ServerName},
		{Name: "health", Method: "GET", Path: "/api/v1/gameserverhealth"},
	}
	router, err = rata.NewRouter(rteRoutes, rteHandlers)
	return
}
