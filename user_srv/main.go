package main

import (
	"user_srv/core"
	"user_srv/initialize"
)

func main() {
	initialize.InitConfig()
	initialize.InitConfigFromNacos()
	initialize.InitLogger()
	initialize.InitSrvConn()

	core.RegisterService()
	core.MainExit()
}
