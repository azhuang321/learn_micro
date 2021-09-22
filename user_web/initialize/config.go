package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop_api/user_web/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	clientConfig := constant.ClientConfig{
		NamespaceId:          "d47a8dad-2d8a-4f2a-a179-186228ead0e9", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:            5000,
		NotLoadCacheAtStart:  true,
		UpdateCacheWhenEmpty: true,
		LogDir:               "logs/nacos/log",
		CacheDir:             "logs/nacos/cache",
		RotateTime:           "1h",
		MaxAge:               3,
		LogLevel:             "error",
		Username:             "nacos",
		Password:             "nacos",
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "127.0.0.1",
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
	}
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		zap.S().Errorf("初始Nacos的Client错误:%s", err.Error())
		return
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "user-web.json",
		Group:  "dev",
	})

	if err != nil {
		zap.S().Errorf("获取Nacos的Client错误:%s", err.Error())
		return
	}
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: "user-web.json",
		Group:  "dev",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})
	if err != nil {
		zap.S().Errorf("监听Nacos的Client错误:%s", err.Error())
		return
	}

	err = json.Unmarshal([]byte(content), global.ServerConfig)
	if err != nil {
		zap.S().Errorf("获取Nacos配置错误:%s", err.Error())
		return
	}
	fmt.Printf("%+v\n", global.ServerConfig)
}

func InitConfig1() {
	data := GetEnvInfo("Debug")
	var configFileName string
	configFileNamePrefix := "config"
	if data {
		configFileName = fmt.Sprintf("config_file/%s_debug.yaml", configFileNamePrefix)
	} else {
		configFileName = fmt.Sprintf("config_file/%s_pro.yaml", configFileNamePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息：%v", *global.ServerConfig)

	go func() {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			zap.S().Infof("配置信息改变：%s", e.Name)
			_ = v.ReadInConfig() // 读取配置数据
			_ = v.Unmarshal(global.ServerConfig)
			zap.S().Infof("配置信息：%v", *global.ServerConfig)
		})
	}()
}
