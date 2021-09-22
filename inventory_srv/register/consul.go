package register

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"inventory_srv/global"
)

type ConsulRegister struct {
	ConsulCent *api.Client
}

func NewConsulRegister() (ConsulRegister, error) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port)
	client, err := api.NewClient(cfg)

	consulRegister := ConsulRegister{}
	if err != nil {
		return consulRegister, err
	}
	consulRegister.ConsulCent = client
	return consulRegister, nil
}

func (c ConsulRegister) Register(name, id, address string, port int, tags []string, check *api.AgentServiceCheck) bool {
	if check == nil {
		//生成对应的检查对象
		check = &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", address, port),
			GRPCUseTLS:                     false,
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		}
	}
	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = address
	registration.Check = check

	err := c.ConsulCent.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Errorf("注册服务中心失败:%s", err.Error())
		return false
	}
	return true
}

func (c ConsulRegister) Deregister(serviceId string) bool {
	err := c.ConsulCent.Agent().ServiceDeregister(serviceId)
	if err != nil {
		zap.S().Errorf("下线服务中心失败:%s", err.Error())
		return false
	}
	return true
}

func (c ConsulRegister) GetAllService() {
	data, err := c.ConsulCent.Agent().Services()
	if err != nil {
		zap.S().Errorf("下线服务中心失败:%s", err.Error())
		return
	}
	for key, _ := range data {
		fmt.Println(key)
	}
}

func (c ConsulRegister) FilterService(filter string) {
	data, err := c.ConsulCent.Agent().ServicesWithFilter(`Service == "user-web"`)
	if err != nil {
		zap.S().Errorf("下线服务中心失败:%s", err.Error())
		return
	}
	for key, _ := range data {
		fmt.Println(key)
	}
}
