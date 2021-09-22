package initialize

import (
	"mxshop_api/user_web/global"
	"mxshop_api/user_web/utils"
)

func InitPort() {
	data := GetEnvInfo("Debug")
	data = true
	//线上使用随机端口号
	if data {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

}
