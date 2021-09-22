package library

import (
	"flag"
	"github.com/gookit/ini/v2"
	"log"
)

func InitConfig() {
	err := ini.LoadExists("config/database.ini")
	if err != nil {
		log.Fatal("加载数据库配置失败")
		return
	}
	err = ini.LoadExists("config/app.ini")
	if err != nil {
		log.Fatal("加载项目配置失败")
		return
	}
	err = ini.LoadExists("config/consul_config.ini")
	if err != nil {
		log.Fatal("加载consul配置失败")
		return
	}
}

func GetArgs() map[string]interface{} {
	var host string
	var port int
	flag.StringVar(&host, "h", "172.17.0.1", "主机名,默认 172.17.0.1")
	flag.IntVar(&port, "p", 0, "端口号,默认为 0")
	flag.Parse()

	if port == 0 {
		port, _ = GetFreePort()
		//port = 8000
	}

	return map[string]interface{}{
		"host": host,
		"port": port,
	}
}
