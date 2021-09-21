package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"user_srv/global"
	"user_srv/handler"
	"user_srv/initialize"
	"user_srv/proto/gen/user_pb"
	"user_srv/register"
	"user_srv/utils"
)

func initFrameWork() {
	initialize.InitConfig()
	initialize.InitLogger()
}

func main() {
	initFrameWork()

	gs := grpc.NewServer()

	//注册用户服务
	userService := &handler.UserService{}
	user_pb.RegisterUserServer(gs, userService)

	//注册健康检查服务
	healthCheckSrv := &handler.HealthCheckSrv{Status: grpc_health_v1.HealthCheckResponse_SERVING, Reason: "running"}
	grpc_health_v1.RegisterHealthServer(gs, healthCheckSrv)

	args := utils.GetArgs()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", args["host"], args["port"]))

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

	u2 := uuid.NewV4()
	zap.S().Info("开始注册服务中心....")
	consulRegister, err := register.NewConsulRegister()
	if err != nil {
		zap.S().Errorf("注册服务中心失败:%s", err.Error())
	} else {
		if consulRegister.Register(global.Config.ProjectName, fmt.Sprintf("%s", u2), args["host"].(string), args["port"].(int), global.Config.Consul.Tags, nil) {
			zap.S().Info("注册服务中心成功")
		}
	}

	//主进程信号退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.S().Info("服务关闭中 ...")
	zap.S().Info("注销服务中心...")
	if consulRegister.Deregister(fmt.Sprintf("%s", u2)) {
		zap.S().Info("注销服务中心成功")
	}
}
