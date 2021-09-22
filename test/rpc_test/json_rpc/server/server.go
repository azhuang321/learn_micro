package main

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloService struct{}

func (s *HelloService) Hello(request string, reply *string) error {
	*reply = "hello " + request
	return nil
}

func main() {
	rpc.RegisterName("HelloService", new(HelloService))
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic("启动错误")
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("接收")
		}
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}

}
