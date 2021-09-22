package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop_api/user_web/api"
	middlewares "mxshop_api/user_web/middlerwares"
)

func InitUserRouter(Router *gin.RouterGroup) {
	zap.S().Infof("配置用户相关router")
	UserRouterGroup := Router.Group("user").Use(middlewares.Trace())
	{
		UserRouterGroup.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouterGroup.POST("pwd_login", api.PasswordLogin)
		UserRouterGroup.POST("register", api.Register)
	}
}
