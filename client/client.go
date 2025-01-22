package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIP string
	ServerPort int
	Name string
	conn net.Conn
	flag int // 当前客户端模式
}

func NewClient(serverIP string, serverPort int) *Client {
	client := &Client{
		ServerIP: serverIP,
		ServerPort: serverPort,
		flag: 999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("net.Dial err: ", err)
		return nil
	}

	client.conn = conn

	return client
}

// 处理服务端响应
func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("---------欢迎使用即时通讯系统---------")
	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("0. 退出系统")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("输入错误，请重新输入数字")
		return false
	}
}

func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println("请输入聊天内容，输入exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err: ", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println("请输入聊天内容，输入exit退出")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var chatMsg string
	var remoteName string

	client.SelectUsers()
	fmt.Println("请输入聊天对象用户名，输入exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("请输入聊天内容，输入exit退出")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write err: ", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println("请输入聊天内容，输入exit退出")
			fmt.Scanln(&chatMsg)
		}
		client.SelectUsers()
		fmt.Println("请输入聊天对象用户名，输入exit退出")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) UpdateName() bool{
	fmt.Println("请输入用户名：")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return false
	}
	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		if client.menu() {
			switch client.flag {
			case 1:
				client.PublicChat()
			case 2:
				client.PrivateChat()
			case 3:
				client.UpdateName()
			default:
			}
		}
	}
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

	go client.DealResponse()

	fmt.Println("NewClient success")

	client.Run()
}