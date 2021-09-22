package initialize

import (
	"github.com/gin-gonic/gin"
	middlewares "mxshop_api/user_web/middlerwares"
	"mxshop_api/user_web/router"
	"net/http"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	Router.Use(middlewares.Cors())
	ApiRouter := Router.Group("/v1")
	router.InitUserRouter(ApiRouter)
	router.InitBaseRouter(ApiRouter)

	return Router
}
