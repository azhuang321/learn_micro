package global

import (
	"mxshop_api/user_web/config"
	"mxshop_api/user_web/proto"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	UserSrvClient proto.UserClient
)
