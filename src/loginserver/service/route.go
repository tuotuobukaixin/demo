package service

import (
	"net/http"

	"github.com/tedsuo/rata"
)

var rteRoutes = rata.Routes{
	{Name: "get_user", Method: "GET", Path: "/api/v1/user"},
	{Name: "add_user", Method: "POST", Path: "/api/v1/user"},
	{Name: "update_user", Method: "PUT", Path: "/api/v1/user"},
}
var rteHandlers = rata.Handlers{
	"get_user":    http.HandlerFunc(GetUser),
	"add_user":    http.HandlerFunc(Adduser),
	"update_user": http.HandlerFunc(Updateuser),
}

//GetRouter return the router fo REST service
// @SubApi k8s-runtime API [/k8sruntime]
func GetRouter() (router http.Handler, err error) {
	router, err = rata.NewRouter(rteRoutes, rteHandlers)
	return
}
