package service

import (
	"net/http"

	"github.com/tedsuo/rata"
)

var rteRoutes = rata.Routes{
	{Name: "get_demotest", Method: "GET", Path: "/api/v1/demotest"},
	{Name: "add_demotest", Method: "POST", Path: "/api/v1/demotest"},
	{Name: "update_demotest", Method: "PUT", Path: "/api/v1/demotest"},
}
var rteHandlers = rata.Handlers{
	"get_demotest":    http.HandlerFunc(GetDemoTest),
	"add_demotest":    http.HandlerFunc(AddDemoTest),
	"update_demotest": http.HandlerFunc(UpdateDemoTest),
}


func GetRouter() (router http.Handler, err error) {
	router, err = rata.NewRouter(rteRoutes, rteHandlers)
	return
}
