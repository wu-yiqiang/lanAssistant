package main

import (
	_ "fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

func NewUser(conn net.Conn) *User {
	/* 创建用户API */
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	go user.ListenMessage()

	return user

}


func (this *User) ListenMessage() {
	/* 监听当前user channel的方法 */
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
