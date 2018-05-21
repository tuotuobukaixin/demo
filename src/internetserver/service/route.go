package service

import (
	"net/http"

	"github.com/tedsuo/rata"
)

var rteRoutes = rata.Routes{
	{Name: "get_gameserver", Method: "GET", Path: "/api/v1/gameserver"},
	{Name: "add_gameserver", Method: "POST", Path: "/api/v1/gameserver"},
	{Name: "update_gameserver", Method: "PUT", Path: "/api/v1/gameserver"},
}
var rteHandlers = rata.Handlers{
	"get_gameserver":    http.HandlerFunc(GetGameserver),
	"add_gameserver":    http.HandlerFunc(AddGameserver),
	"update_gameserver": http.HandlerFunc(UpdateGameserver),
}

//GetRouter return the router fo REST service
// @SubApi k8s-runtime API [/k8sruntime]
func GetRouter() (router http.Handler, err error) {
	router, err = rata.NewRouter(rteRoutes, rteHandlers)
	return
}
