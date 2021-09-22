package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.uber.org/zap"
	"user_srv/utils"

	"user_srv/config"
)

func InitConfigFromNacos() {
	clientConfig := constant.ClientConfig{
		NamespaceId:          config.Config.Nacos.NamespaceId, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:            5000,
		NotLoadCacheAtStart:  true,
		UpdateCacheWhenEmpty: true,
		LogDir:               config.Config.Nacos.LogDir,
		CacheDir:             config.Config.Nacos.CacheDir,
		RotateTime:           config.Config.Nacos.RotateTime,
		MaxAge:               config.Config.Nacos.MaxAge,
		LogLevel:             config.Config.Nacos.LogLevel,
		Username:             config.Config.Nacos.Username,
		Password:             config.Config.Nacos.Password,
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      config.Config.Nacos.IpAddr,
			ContextPath: config.Config.Nacos.ContextPath,
			Port:        config.Config.Nacos.Port,
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
		DataId: config.Config.Nacos.DataId,
		Group:  config.Config.Nacos.Group,
	})
	if err != nil {
		zap.S().Errorf("获取Nacos的Client错误:%s", err.Error())
		return
	}
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: config.Config.Nacos.DataId,
		Group:  config.Config.Nacos.Group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})
	if err != nil {
		zap.S().Errorf("监听Nacos的Client错误:%s", err.Error())
		return
	}

	err = json.Unmarshal([]byte(content), config.Config)
	if err != nil {
		zap.S().Errorf("获取Nacos配置错误:%s", err.Error())
		return
	}
	fmt.Println("Nacos配置信息：")
	utils.PrettyPrint(*config.Config)
}
