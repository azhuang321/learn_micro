package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"mxshop_api/user_web/global"
	"mxshop_api/user_web/initialize"
	"mxshop_api/user_web/register"
)

//初始化框架
func initFrame() {
	initialize.InitLogger("debug", "logs/test.log", 1, 5, 1, false)
	initialize.InitConfig()
	initialize.InitTranslate()
	initialize.InitSrvConn()
	initialize.InitPort()
	initialize.InitSentinel()
}

func main() {
	initFrame()

	u2 := uuid.NewV4()
	zap.S().Info("开始注册服务中心....")
	consulRegister, err := register.NewConsulRegister()
	if err != nil {
		zap.S().Errorf("注册服务中心失败:%s", err.Error())
	} else {
		if consulRegister.Register(global.ServerConfig.Name, fmt.Sprintf("%s", u2), "172.17.0.1", global.ServerConfig.Port, global.ServerConfig.Tags, nil) {
			zap.S().Info("注册服务中心成功")
		}
	}

	// 初始化router
	Router := initialize.Routers()
	go func() {
		zap.S().Infof("启动web服务，端口：%d", global.ServerConfig.Port)
		if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动web服务失败：", err.Error())
		}
	}()

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
