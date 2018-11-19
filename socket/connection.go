package socket

import (
	"github.com/gorilla/websocket"
	"sync"
	"errors"
)

type Connection struct {
	wsConn *websocket.Conn
	inChan chan []byte
	outChan chan []byte
	closeChan chan byte
	mutex sync.Mutex
	isClosed bool
}

func initConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn: wsConn,
		inChan: make(chan []byte, 1000),
		outChan: make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
	}

	// 读协程
	go conn.readLoop()

	// 写协程
	go conn.writeLoop()

	return
}

// 读取消息
func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <- conn.inChan:
	case <- conn.closeChan:
		err = errors.New("connection is close")
	}
	return
}

// 写入消息
func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <- conn.closeChan:
		err = errors.New("connection is close")
	}
	return
}

// 关闭连接
func (conn *Connection) Close()  {
	// 线程安全的 可重入的Close
	conn.wsConn.Close()
	// 加锁
	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

func (conn *Connection) readLoop() {
	var (
		data []byte
		err error
	)
	for  {
		if _, data, err = conn.wsConn.ReadMessage(); err!= nil {
			// 读取失败 跳出循环 关闭连接
			goto ERR
		}
		select {
		case conn.inChan <- data:
		case <- conn.closeChan:
			// closeChan 被关闭时调用
			goto ERR
		}
	}
	// 标签
	ERR:
		conn.Close()
}

func (conn *Connection) writeLoop() {
	var (
		data []byte
		err error
	)
	for  {
		select {
		case data = <- conn.outChan:
		case <- conn.closeChan:
			goto ERR
		}

		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			// 发送失败 跳出循环 关闭连接
			goto ERR
		}
	}
	// 标签
	ERR:
		conn.Close()
}