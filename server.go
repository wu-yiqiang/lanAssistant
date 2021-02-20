package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	//创建一个server接口
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

func (this *Server) BroadCast(user *User, msg string) {
	/* 广播消息 */
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) ListenChannel() {
	/* 监听channel */
	for {
		msg := <-this.Message
		//发送消息
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) Handler(conn net.Conn) {
	/*当前连接的业务 */
	user := NewUser(conn)
	//将上线用户进行广播
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	//广播用户上线消息
	this.BroadCast(user, "已上线")
	//接收客户端发送的消息
	go func() {
	    buf :=make ([]byte,4096)
		for{
			n,err:=conn.Read(buf)
			if n==0 {
				this.BroadCast(user,"下线")
				return 
			}
			if err!=nil&&err !=io.EOF {
				fmt.Println("conn Read err:",err)
			}
			//提取用户消息
			msg :=string(buf[:n-1])	
			
			//将得到的消息进行广播
			this.BroadCast(user,msg)
			
		}
	}
	select {}
}

func (this *Server) Start() {
	/* 启动服务器的接口 */

	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net listen err:", err)
		return
	}
	defer listener.Close()
	//监听go协程
	go this.ListenChannel()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("litstener accept err:", err)
			continue
		}



		go this.Handler(conn)
	}

}
