package service

import (
	"net/http"
	"websocket/http/controller"
	"websocket/socket"
)

// 路由
func route()  {
	// websockets 路由
	http.HandleFunc("/wss", socket.Websocket)

	// http 根路由
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		home := new(controller.HomeController)
		home.Writer = writer
		home.Request = request
		home.Index()
	})
}