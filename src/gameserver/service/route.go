package service

import (
	"net/http"

	"github.com/tedsuo/rata"
)

var rteRoutes = rata.Routes{
	{Name: "get_gameserver", Method: "GET", Path: "/api/v1/gameserverdetail"},
	{Name: "health", Method: "GET", Path: "/api/v1/gameserverhealth"},
}
var rteHandlers = rata.Handlers{
	"get_gameserver":    http.HandlerFunc(GetGameserverDetail),
	"health":    http.HandlerFunc(Health),
}

//GetRouter return the router fo REST service
// @SubApi k8s-runtime API [/k8sruntime]
func GetRouter() (router http.Handler, err error) {
	router, err = rata.NewRouter(rteRoutes, rteHandlers)
	return
}
