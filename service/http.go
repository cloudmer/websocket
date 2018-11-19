package service

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 10,
	WriteBufferSize: 1024 * 10,
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// http 协议 升级为 websocket
func wss(writer http.ResponseWriter, request *http.Request)  {
	fmt.Println(request.Cookies())
	// 设置websocket 下的cookie
	for i := range request.Cookies() {
		cookieIndex := request.Cookies()[i]
		fmt.Println(cookieIndex)
		cookie := http.Cookie{Name: cookieIndex.Name, Value: cookieIndex.Value, Expires: cookieIndex.Expires}
		http.SetCookie(writer, &cookie)
	}
	conn, err := upgrader.Upgrade(writer, request, writer.Header())
	if err != nil {
		// http 协议 升级 为 websocket 协议失败
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		fmt.Println(string(msg))
		if err == nil {
			conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

// 启动服务
func Start()  {
	// 路由解析
	route()
	// 启动服务
	//http.ListenAndServeTLS("0.0.0.0:10086", "/Users/cloud/wwwroot/golang/src/websocket/secret/server.crt", "/Users/cloud/wwwroot/golang/src/websocket/secret/server.key", nil)
	http.ListenAndServe("0.0.0.0:10086", nil)
}