package main

import "net"

type User struct {
	Name string
	Addr string
	C chan string
	conn net.Conn
	server *Server
}

func NewUser(conn net.Conn,server *Server) *User{
	user:=&User{
		Name: conn.RemoteAddr().String(),
		Addr: conn.RemoteAddr().String(),
		C:    make(chan string),
		conn: conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}

func (this *User) Online(){
	// 用户上线 将用户加入 OnlineMap
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	// 广播用户上线信息
	this.server.BoardCast(this,"已上线")
}

func (this *User) OffLine(){
	// 用户下线 将用户从 OnlineMap 删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap,this.Name)
	this.server.mapLock.Unlock()
	// 广播用户下线信息
	this.server.BoardCast(this,"下线")
}

// 用户消息处理
func (this *User) DoMessage(msg string){
	this.server.BoardCast(this,msg)
}

func (this *User) ListenMessage(){
	for  {
		msg:=<-this.C
		this.conn.Write([]byte(msg+"\n"))
	}
}