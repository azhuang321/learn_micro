package initialize

import (
	"coupons_srv/utils"
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.uber.org/zap"

	. "coupons_srv/config"
)

func InitConfigFromNacos() {
	clientConfig := constant.ClientConfig{
		NamespaceId:          Config.Nacos.NamespaceId, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:            5000,
		NotLoadCacheAtStart:  true,
		UpdateCacheWhenEmpty: true,
		LogDir:               Config.Nacos.LogDir,
		CacheDir:             Config.Nacos.CacheDir,
		RotateTime:           Config.Nacos.RotateTime,
		MaxAge:               Config.Nacos.MaxAge,
		LogLevel:             Config.Nacos.LogLevel,
		Username:             Config.Nacos.Username,
		Password:             Config.Nacos.Password,
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      Config.Nacos.IpAddr,
			ContextPath: Config.Nacos.ContextPath,
			Port:        Config.Nacos.Port,
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
		DataId: Config.Nacos.DataId,
		Group:  Config.Nacos.Group,
	})
	if err != nil {
		zap.S().Errorf("获取Nacos的Client错误:%s", err.Error())
		return
	}
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: Config.Nacos.DataId,
		Group:  Config.Nacos.Group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})
	if err != nil {
		zap.S().Errorf("监听Nacos的Client错误:%s", err.Error())
		return
	}

	err = json.Unmarshal([]byte(content), Config)
	if err != nil {
		zap.S().Errorf("获取Nacos配置错误:%s", err.Error())
		return
	}
	fmt.Println("Nacos配置信息：")
	utils.PrettyPrint(*Config)
}
