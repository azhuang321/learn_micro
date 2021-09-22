package main

import (
	"user_srv/core"
	"user_srv/initialize"
)

func main() {
	initialize.InitConfig()
	initialize.InitConfigFromNacos()
	initialize.InitLogger()

	core.RegisterService()
	core.MainExit()
}
