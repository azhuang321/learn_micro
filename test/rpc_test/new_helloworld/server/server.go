package main

import (
	"myrpc/new_helloworld/handler"
	"myrpc/new_helloworld/server_proxy"
	"net"
	"net/rpc"
)

func main() {
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic("监听端口失败")
	}

	server_proxy.RegisterHelloService(&handler.HelloService{})

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("建立连接失败")
		}
		go rpc.ServeConn(conn)
	}
}
