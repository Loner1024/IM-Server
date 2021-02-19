package main

import (
	"fmt"
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

// 给当前用户发送消息
func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg))
}

// 用户消息处理
func (this *User) DoMessage(msg string) {
	fmt.Println(len(msg), msg[:3])
	if msg == "who" {
		for _, user := range this.server.OnlineMap {
			this.SendMessage("[" + user.Addr + "]" + user.Name + "在线...\n")
		}
		return
	}
	if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		if _, ok := this.server.OnlineMap[newName]; !ok {
			// 修改 name
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.SendMessage("名字已修改为: " + newName + "\n")
			return
		} else {
			this.SendMessage("名字已经存在\n")
		}
	}
	if len(msg) > 4 && msg[:3] == "to|" {
		fmt.Println("here")
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMessage("消息格式不正确\n")
			return
		}
		// 找到 User 对象
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMessage("用户不存在\n")
			return
		}
		// 发送消息出去
		remoteUser.SendMessage(this.Name + "对你说 " + strings.Split(msg, "|")[2] + "\n")
		return
	}
	this.server.BoardCast(this, msg)
}

func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
