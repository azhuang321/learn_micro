package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"

	"google.golang.org/grpc"

	"mygrpc/proto"
)

type Server struct {
}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Println("get metadata error")
	}
	for key, val := range md {
		fmt.Println(key, val)
	}
	if usernameSlice, ok := md["username"]; ok {
		fmt.Println(usernameSlice)
	}
	return &proto.HelloReply{
		Message: "Hello " + request.Name,
		Hobby:   request.Hobby,
		Sex:     request.Sex,
		Mp:      request.Mp,
	}, nil
}

func main() {
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("收到新的请求")
		resp, err = handler(ctx, req)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return resp, status.Error(codes.Unauthenticated, "无token")
		}
		var token string

		if val, ok := md["token"]; ok {
			token = val[0]
		}

		if token != "test1" {
			return resp, status.Error(codes.Unauthenticated, "token验证失败")
		}

		resp, err = handler(ctx, req)
		fmt.Println("请求结束")
		return resp, err
	}
	opt := grpc.UnaryInterceptor(interceptor)
	g := grpc.NewServer(opt)
	s := Server{}
	proto.RegisterGreeterServer(g, &s)
	lis, err := net.Listen("tcp", fmt.Sprintf(":8000"))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}
	g.Serve(lis)
}
