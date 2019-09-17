package main

import (
	"fmt"
  "io"
	"net"
	"os"
	"strconv"
)


var	ip = "0.0.0.0"
var	port = 18080
var	s5port = 28080


func main(){
	args := os.Args
	port,_ = strconv.Atoi(args[1])
	s5port,_ = strconv.Atoi(args[2])

	lis, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ip), port, ""})
	ErrHandler(err)
	defer lis.Close()
	fmt.Println("Listen :",port)
	s5lis, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ip), s5port, ""})
	ErrHandler(err)
	defer s5lis.Close()
	fmt.Println("socks5 port:",s5port)

	Server(lis,s5lis)
}


func Server(listen *net.TCPListener,s5listen *net.TCPListener) {
    for {
        s5conn, err := s5listen.Accept()
        if err != nil {
            fmt.Println("接受客户端连接异常:", err.Error())
            continue
        }
        fmt.Println("用户客户端连接来自:", s5conn.RemoteAddr().String())
        defer s5conn.Close()


        conn, err := listen.Accept()
        if err != nil {
            fmt.Println("接受客户端连接异常:", err.Error())
            continue
        }
        fmt.Println("肉鸡客户端连接来自:", conn.RemoteAddr().String())
        defer conn.Close()

        go handle(conn, s5conn)
    }
}


func handle(sconn net.Conn, dconn net.Conn) {
	defer sconn.Close()
  defer dconn.Close()
	ExitChan := make(chan bool, 1)
	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		io.Copy(dconn, sconn)
		ExitChan <- true
	}(sconn, dconn, ExitChan)

	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		io.Copy(sconn, dconn)
		ExitChan <- true
	}(sconn, dconn, ExitChan)
	<-ExitChan
	dconn.Close()
}

func ErrHandler(err error) {
	if err != nil {
		panic(err)
	}
}
