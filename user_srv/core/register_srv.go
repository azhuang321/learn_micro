package core

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	. "user_srv/config"
	"user_srv/handler"
	"user_srv/proto/gen/user_pb"
	"user_srv/register"
	"user_srv/utils"
)

func RegisterService() {
	var err error
	register.SrvRegister, err = register.NewConsulRegister()
	if err != nil {
		zap.S().Errorf("注册服务中心失败:%s", err.Error())
	}

	gs := grpc.NewServer()
	//注册用户服务
	userService := &handler.UserService{}
	user_pb.RegisterUserServer(gs, userService)

	//注册健康检查服务
	healthCheckSrv := &handler.HealthCheckSrv{Status: grpc_health_v1.HealthCheckResponse_SERVING, Reason: "running"}
	grpc_health_v1.RegisterHealthServer(gs, healthCheckSrv)

	args := utils.GetArgs()
	//lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", args["host"], args["port"]))
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 9000))

	if err != nil {
		zap.S().Fatalf("监听端口出现错误:%s", err.Error())
		return
	}

	go func() {
		zap.S().Infof("启动服务成功:%s:%d", args["host"], args["port"])
		if err := gs.Serve(lis); err != nil {
			zap.S().Errorf("启动服务失败:%s\n", err.Error())
		}
	}()

	zap.S().Info("开始注册服务中心....")
	if register.SrvRegister.Register(Config.ProjectName, Config.Uuid, args["host"].(string), args["port"].(int), Config.Consul.Tags, nil) {
		zap.S().Info("注册服务中心成功")
	}
}
