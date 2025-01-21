package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIP string
	ServerPort int
	Name string
	conn net.Conn
}

func NewClient(serverIP string, serverPort int) *Client {
	client := &Client{
		ServerIP: serverIP,
		ServerPort: serverPort,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("net.Dial err: ", err)
		return nil
	}

	client.conn = conn

	return client
}

var serverIP string
var serverPort int

// 命令行参数：client -ip 127.0.0.1 -port 8080
func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "设置 server ip")
	flag.IntVar(&serverPort, "port", 8080, "设置 server port")
}	

func main() {
	// 解析命令行参数
	flag.Parse()

	client := NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println("NewClient err")
		return
	}

	fmt.Println("NewClient success")

	select {}
}