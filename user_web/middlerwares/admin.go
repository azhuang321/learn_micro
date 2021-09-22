package middlewares

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/user_web/models"
	"net/http"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		userInfo := claims.(*models.CustomClaims)

		if userInfo.AuthorityId != 2 {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
