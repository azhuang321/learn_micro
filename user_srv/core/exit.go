package core

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	. "user_srv/config"
	"user_srv/register"
)

var QuitSignal = make(chan os.Signal)

func MainExit() {
	//主进程信号退出
	signal.Notify(QuitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-QuitSignal
	zap.S().Info("服务关闭中 ...")
	zap.S().Info("注销服务中心...")
	if register.SrvRegister.Deregister(Config.Uuid) {
		zap.S().Info("注销服务中心成功")
	}
}
