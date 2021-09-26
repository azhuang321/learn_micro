package global

import (
	"mxshop_api/config"
	"mxshop_api/proto"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	UserSrvClient proto.UserClient
)
