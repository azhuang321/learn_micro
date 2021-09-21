package initialize

import (
	"fmt"
	"user_srv/utils"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"user_srv/global"
)

// GetEnvInfo 读取系统环境变量
func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}

func InitConfig() {
	runMod := GetEnvInfo("PRO")
	var configFileName string
	configFileNamePrefix := "config"
	if runMod == "PRO" {
		runMod = "pro"
		configFileName = fmt.Sprintf("%s_pro.yaml", configFileNamePrefix)
	} else {
		runMod = "debug"
		configFileName = fmt.Sprintf("%s_debug.yaml", configFileNamePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if err := v.Unmarshal(global.Config); err != nil {
		panic(err)
	}
	global.Config.RunMod = runMod
	fmt.Println("配置信息：")
	utils.PrettyPrint(*global.Config)

	//动态监控配置文件变化
	go func() {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			fmt.Printf("配置信息改变：%s", e.Name)
			_ = v.ReadInConfig() // 读取配置数据
			_ = v.Unmarshal(global.Config)
			global.Config.RunMod = runMod
			fmt.Printf("配置信息：%+v\n", *global.Config)
		})
	}()
}
