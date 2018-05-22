package service

import (
	"net/http"

	"github.com/tedsuo/rata"
)

var rteRoutes = rata.Routes{
	{Name: "get_info", Method: "GET", Path: "/api/v1/info"},
}
var rteHandlers = rata.Handlers{
	"get_info":    http.HandlerFunc(GetInfo),
}

//GetRouter return the router fo REST service
// @SubApi k8s-runtime API [/k8sruntime]
func GetRouter() (router http.Handler, err error) {
	router, err = rata.NewRouter(rteRoutes, rteHandlers)
	return
}
