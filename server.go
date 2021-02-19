package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type Server struct {
	IP string
	Port int
	OnlineMap map[string]*User
	mapLock sync.RWMutex
	
	Message chan string
}

func NewServer(ip string,port int) *Server{
	return &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

func (this *Server) ListenMessage(){
	for{
		msg:=<-this.Message
		this.mapLock.Lock()
		for _,cli:=range this.OnlineMap{
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}


func (this *Server) BoardCast(user *User,msg string){
	sendMsg:="["+user.Addr+"]"+user.Name+":"+msg
	this.Message<-sendMsg
}


func (this *Server) Handler(conn net.Conn){
	user:=NewUser(conn,this)
	
	// 用户上线 将用户加入 OnlineMap
	user.Online()
	
	// 接收客户端发送的消息
	go func() {
		buf:=make([]byte,4096)
		for  {
			n,err:=conn.Read(buf)
			if n==0{
				user.OffLine()
				return
			}
			if err!=nil && err!= io.EOF{
				fmt.Println("Conn Read err:",err)
				return
			}
			// 提取用户信息
			msg:=string(buf[:n-1])
			// 广播消息
			user.DoMessage(msg)
		}
	}()
	
	// 阻塞当前 Handler
	select {}
	
}


// 启动服务器的接口
func (this *Server) Start(){
	// socket listen
	listener,err:=net.Listen("tcp",fmt.Sprintf("%s:%d",this.IP,this.Port))
	if err!=nil{
		log.Println("net.Listen err:",err)
		return
	}
	// close listen socket
	defer listener.Close()
	
	// 启动监听 Message 的 goroutine
	go this.ListenMessage()
	
	for {
		// accept
		conn,err:=listener.Accept()
		if err!=nil{
			log.Println("Listener accept err:",err)
			continue
		}
		go this.Handler(conn)
	}
}