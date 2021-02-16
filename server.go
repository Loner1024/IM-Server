package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	IP string
	Port int
}

func NewServer(ip string,port int) *Server{
	return &Server{
		IP:   ip,
		Port: port,
	}
}


func Handler(conn net.Conn){
	log.Println("连接建立")
}

func Start(host *Server){
	listener,err:=net.Listen("tcp",fmt.Sprintf("%s:%d",host.IP,host.Port))
	if err!=nil{
		log.Println("net.Listen err:",err)
		return
	}
	defer listener.Close()
	for {
		conn,err:=listener.Accept()
		if err!=nil{
			log.Println("Listener accept err:",err)
			continue
		}
		go Handler(conn)
	}
}