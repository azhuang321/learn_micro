package initialize

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"

	. "user_srv/config"
)

// GetEnvInfo 读取系统环境变量
func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}

func InitConfig() {
	Config.Uuid = uuid.NewV4().String()

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

	if err := v.Unmarshal(Config); err != nil {
		panic(err)
	}
	Config.RunMod = runMod

	//动态监控配置文件变化
	go func() {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			fmt.Printf("配置信息改变：%s", e.Name)
			_ = v.ReadInConfig() // 读取配置数据
			_ = v.Unmarshal(Config)
			Config.RunMod = runMod
			fmt.Printf("配置信息：%+v\n", *Config)
		})
	}()
}
