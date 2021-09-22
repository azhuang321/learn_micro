package register

import "github.com/hashicorp/consul/api"

type Register interface {
	Register(name, id, address string, port int, tags []string, check *api.AgentServiceCheck) bool
	Deregister(serviceId string) bool
	GetAllService()
	FilterService(filter string)
}
