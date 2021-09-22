package main

import (
	"fmt"
	"inventory_srv/global"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"inventory_srv/handler"
	"inventory_srv/initialize"
	"inventory_srv/library"
	"inventory_srv/proto"
	"inventory_srv/register"
)

func main() {
	library.InitConfig()
	library.InitLogger("debug", "logs/test.log", 1, 5, 1, false)
	initialize.InitConfig()
	args := library.GetArgs()

	g := grpc.NewServer()
	//注册用户服务
	userSrv := &handler.InventoryService{}
	proto.RegisterInventoryServer(g, userSrv)

	//注册健康检查服务
	healthCheckSrv := &handler.HealthCheckSrv{Status: grpc_health_v1.HealthCheckResponse_SERVING, Reason: "running"}
	grpc_health_v1.RegisterHealthServer(g, healthCheckSrv)

	//lis,err := net.Listen("tcp", fmt.Sprintf("%s:%d",args["host"],args["port"]))
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 8000))
	if err != nil {
		zap.S().Errorf("启动服务失败:%s\n", err.Error())
		return
	}

	go func() {
		zap.S().Infof("启动服务成功:%s:%d", args["host"], args["port"])
		if err := g.Serve(lis); err != nil {
			zap.S().Errorf("启动服务失败:%s\n", err.Error())
		}
	}()

	u2 := uuid.NewV4()
	zap.S().Info("开始注册服务中心....")
	consulRegister, err := register.NewConsulRegister()
	if err != nil {
		zap.S().Errorf("注册服务中心失败:%s", err.Error())
	} else {
		if consulRegister.Register(global.ServerConfig.Name, fmt.Sprintf("%s", u2), args["host"].(string), args["port"].(int), global.ServerConfig.Tags, nil) {
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
