package socket

import (
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// http协议 升级为 websocket 协议
func Websocket(writer http.ResponseWriter, request *http.Request)  {
	var (
		wsConn *websocket.Conn
		conn *Connect
		data []byte
		err error
	)
	// http协议 升级为 websocket 协议
	if wsConn, err = upgrader.Upgrade(writer, request, nil); err != nil {
		// 协议升级失败
		return
	}

	// 新的连接加入
	if conn, err = initConnect(wsConn); err != nil {
		goto ERR
	}

	for  {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}

	// 标签
	ERR:
		//TODO 关闭连接操作
		wsConn.Close()
}