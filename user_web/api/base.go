package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误" + e.Message(),
				})
			}
		}
	}
	return
}

// 错误快捷返回
func errReturn(c *gin.Context, errCode int, errMsg string) {
	c.JSON(errCode, gin.H{
		"msg": errMsg,
	})
}

// 绑定参数并验证
func bindJsonAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		errReturn(c, http.StatusBadRequest, "参数错误")
		return false
	}
	v := validate.Struct(obj)
	if !v.Validate() {
		errReturn(c, http.StatusBadRequest, v.Errors.One())
		return false
	}
	return true
}
