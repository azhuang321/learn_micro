package main

import (
	"coupons_srv/core"
	"coupons_srv/initialize"
)

func main() {
	initialize.InitConfig()
	initialize.InitConfigFromNacos()
	initialize.InitLogger()

	core.RegisterService()
	core.MainExit()
}
