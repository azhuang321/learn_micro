package main

import (
	"net"

	"google.golang.org/grpc"

	"user_srv/handler"
	"user_srv/proto/gen/user_pb"
)

func main() {
	gs := grpc.NewServer()

	userService := handler.UserService{}
	user_pb.RegisterUserServer(gs, userService)

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	err = gs.Serve(lis)
	if err != nil {
		panic(err)
	}
}
