package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.uber.org/zap"
	"inventory_srv/global"
)

func InitConfig() {
	clientConfig := constant.ClientConfig{
		NamespaceId:          "74ac63ef-a03a-40e6-b12a-7ca1af81679a", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
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
		DataId: "inventory-srv.json",
		Group:  "dev",
	})
	if err != nil {
		zap.S().Errorf("获取Nacos的Client错误:%s", err.Error())
		return
	}
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: "inventory-srv.json",
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
}
