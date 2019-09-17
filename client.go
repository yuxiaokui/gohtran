package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)


type sockIP struct {
	A, B, C, D byte
	PORT       uint16
}


func (ip sockIP) toAddr() string {
	return fmt.Sprintf("%d.%d.%d.%d:%d", ip.A, ip.B, ip.C, ip.D, ip.PORT)
}

func handleClientRequest(client net.Conn) {
    if client == nil {
        return
    }
    defer client.Close()

    var b [1024]byte
    n, err := client.Read(b[:])
    if err != nil {
        return
    }

    if b[0] == 0x05 { //只处理Socks5协议
        client.Write([]byte{0x05, 0x00})
        n, err = client.Read(b[:])
        var addr string
        switch b[3] {
        case 0x01:
        		sip := sockIP{}
        		if err := binary.Read(bytes.NewReader(b[4:n]), binary.BigEndian, &sip); err != nil {
        			log.Println("请求解析错误")
        			return
        		}
        		addr = sip.toAddr()
        	case 0x03:
        		host := string(b[5 : n-2])
        		var port uint16
        		err = binary.Read(bytes.NewReader(b[n-2:n]), binary.BigEndian, &port)
        		if err != nil {
        			log.Println(err)
        			return
        		}
        		addr = fmt.Sprintf("%s:%d", host, port)
        	}

        server, err := net.Dial("tcp",addr)
        if err != nil {
            return
        }
        defer server.Close()
        client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功
        //进行转发
        go io.Copy(server, client)
        io.Copy(client, server)
    }

}

func slave(listenTarget string) {
	var RemoteConn net.Conn
	var err error
	for {
    for {
		  RemoteConn, err = net.Dial("tcp", listenTarget)
		  if err == nil {
				break
			}
    }
    go handleClientRequest(RemoteConn)
	}
}



func main() {
	args := os.Args
	slave(args[1])
}
