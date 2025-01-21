package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

const (
	BufSize = 1024
)

// Server 服务器结构体,用于管理在线用户和消息广播
type Server struct {
	IP   string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

// NewServer 创建服务器
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// ListenMessage 监听Message消息
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// BroadCast 广播消息
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

// Handler 处理用户请求
func (this *Server) Handler(conn net.Conn) {

	user := NewUser(conn)

	// 将用户添加到在线用户列表
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// 广播上线消息
	this.BroadCast(user, "用户已上线")

	// 处理用户发送的消息
	go func() {
		buf := make([]byte, BufSize)
		for {
			n, err := conn.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Println("conn.Read err:", err)
				return
			}
			if n == 0 {
				this.BroadCast(user, "用户已下线")
				return
			}
			// 提取用户消息（去掉\n）
			msg := string(buf[:n-1])
			// 将读到的消息广播给其他用户
			this.BroadCast(user, msg)
		}
	}()

	// 阻塞等待用户退出
	select {}
}

// Start 启动服务器
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.IP, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	// close listen socket
	defer listener.Close()

	// listen Message 通道
	go this.ListenMessage()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}
		// do handler
		go this.Handler(conn)
	}
}
