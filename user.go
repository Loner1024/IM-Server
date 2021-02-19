package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	user := &User{
		Name:   conn.RemoteAddr().String(),
		Addr:   conn.RemoteAddr().String(),
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}

func (this *User) Online() {
	// 用户上线 将用户加入 OnlineMap
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	// 广播用户上线信息
	this.server.BoardCast(this, "已上线")
}

func (this *User) OffLine() {
	// 用户下线 将用户从 OnlineMap 删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	// 广播用户下线信息
	this.server.BoardCast(this, "下线")
}

// 给当前用户发送消息（自己给自己发消息
func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg))
}

// 用户消息处理
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		for _, user := range this.server.OnlineMap {
			this.SendMessage("[" + user.Addr + "]" + user.Name + "在线...\n")
		}
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		if _, ok := this.server.OnlineMap[newName]; !ok {
			// 修改 name
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.SendMessage("名字已修改为: " + newName + "\n")
		} else {
			this.SendMessage("名字已经存在\n")
		}
	} else {
		this.server.BoardCast(this, msg)
	}
}

func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
