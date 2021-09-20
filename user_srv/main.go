package main

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"user_srv/handler"
	"user_srv/initialize"
	"user_srv/proto/gen/user_pb"
)

func initFrameWork() {
	initialize.InitLogger("debug", "logs/test.log", 1, 5, 1, false)
}

func main() {
	initFrameWork()

	gs := grpc.NewServer()

	userService := handler.UserService{}
	user_pb.RegisterUserServer(gs, userService)

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		zap.S().Fatalf("监听端口出现错误:%s", err.Error())
		return
	}

	err = gs.Serve(lis)
	if err != nil {
		zap.S().Fatalf("启动服务出现错误:%s", err.Error())
		return
	}
}
