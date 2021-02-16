package main

import (
	"fmt"
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
	user:=NewUser(conn)
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()
	this.BoardCast(user,"已上线")
	
}

func (this *Server) Start(){
	listener,err:=net.Listen("tcp",fmt.Sprintf("%s:%d",this.IP,this.Port))
	if err!=nil{
		log.Println("net.Listen err:",err)
		return
	}
	defer listener.Close()
	
	go this.ListenMessage()
	
	for {
		conn,err:=listener.Accept()
		if err!=nil{
			log.Println("Listener accept err:",err)
			continue
		}
		go this.Handler(conn)
	}
}