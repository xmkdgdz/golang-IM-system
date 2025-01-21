package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// NewUser 创建一个用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}
	go user.ListenMessage()
	return user
}

// Online 用户上线业务
func (this *User) Online() {
	// 将用户添加到在线用户列表
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	// 广播上线消息
	this.server.BroadCast(this, "已上线")
}

// Offline 用户下线业务
func (this *User) Offline() {
		this.server.mapLock.Lock()
		delete(this.server.OnlineMap, this.Name)
		this.server.mapLock.Unlock()
		this.server.BroadCast(this, "已下线")
}

// SendMsg 给当前用户客户端发送消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// DoMessage 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询在线用户
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ": 在线\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|"{
		// 改名: rename|newName
		newName := strings.Split(msg, "|")[1]
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("用户名已存在\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.SendMsg("用户名改为：" + this.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 私聊: to|name|msg
		toName := strings.Split(msg, "|")[1]
		if toName == "" {
			this.SendMsg("私聊对象不能为空\n")
			return
		}
		toUser, ok := this.server.OnlineMap[toName]
		if !ok {
			this.SendMsg("该用户不存在\n")
			return
		}
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("私聊内容不能为空\n")
			return
		}
		toUser.SendMsg("[" + this.Addr + "]" + this.Name + " 对您说: " + content)
	} else {
		this.server.BroadCast(this, msg)
	}
}

// ListenMessage 监听当前 User 通道，一旦有消息，就直接发送给对应客户端
func (this *User) ListenMessage() {
	for msg := range this.C {
		if _, err := this.conn.Write([]byte(msg + "\n")); err != nil {
			break
		}
	}
}
