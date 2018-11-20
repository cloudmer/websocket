package socket

import (
	"github.com/gorilla/websocket"
	"sync"
	"errors"
)

// websocket 连接 接口
type WsConnInterface interface {
	// 读取 消息队列
	ReadMessage() (data []byte, err error)
	// 写入 发送队列
	WriteMessage(data []byte) (err error)
	// 关闭 websocket 连接
	Close()
	// 读取 websocket 发来的消息 加入 readChan 队列
	readLoop()
	// 读取 sendChand 队列消息 发送给 Client 消息
	writeLoop()
}

// websocket 连接对象
type Connect struct {
	wsConn    *websocket.Conn // websocket 连接
	readChan  chan []byte	  // 读取客户端发来的消息 队列 chan
	sendChan  chan []byte	  // 写入需要发送到连接的 队列 chan
	stopChan  chan byte		  // 为保证线程安全 连接关闭时 chan
	stopMutex sync.Mutex	  // isClose 互斥锁
	isStop    bool
}

// 初始化 websocket 连接
func initConnect(wsConn *websocket.Conn) (conn *Connect, err error) {
	conn = &Connect{
		wsConn:   wsConn,
		readChan: make(chan []byte, 1000),
		sendChan: make(chan []byte, 1000),
		stopChan: make(chan byte),
	}

	// 读协程
	go conn.readLoop()
	// 写协程
	go conn.writeLoop()

	return
}

// 读取消息队列
func (conn *Connect) ReadMessage() (data []byte, err error) {
	select {
	case data = <- conn.readChan:
	case <- conn.stopChan:
		errors.New("connect is closed")
	}
	return
}

// 将消息 发送到发送队列
func (conn *Connect) WriteMessage(data []byte) (err error) {
	select {
	case conn.sendChan <- data:
	case <- conn.stopChan:
		errors.New("connect is closed")
	}
	return
}

// 关闭 websocket 连接
func (conn *Connect) Close() {
	// 线程安全的 可重入的Close
	conn.wsConn.Close()
	// 加锁
	conn.stopMutex.Lock()
	if !conn.isStop {
		// 关闭 chan
		close(conn.stopChan)
		conn.isStop = true
	}
	// 解锁
	conn.stopMutex.Unlock()
}

// 读取 websocket 发来的消息 加入 readChan 队列
func (conn *Connect) readLoop() {
	var data []byte
	var err error
	for  {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			// 读取失败 跳出循环 关闭连接
			goto ERR
		}
		select {
		case conn.readChan <- data:
		case <-conn.stopChan:
			goto ERR
		}
	}
	// 标签
	ERR:
		// 关闭连接
		conn.Close()
}

// 读取 sendChan 队列 将消息发送给 客户端
func (conn *Connect) writeLoop() {
	var data []byte
	var err error
	for  {
		select {
		case data = <- conn.sendChan:
		case <- conn.stopChan:
			goto ERR
		}
		// 给客户端发消息
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			// 发送失败 跳出循环 关闭连接
			goto ERR
		}
	}
	ERR:
		conn.Close()
}