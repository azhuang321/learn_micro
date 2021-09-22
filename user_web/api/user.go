package api

import (
	"context"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"mxshop_api/user_web/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop_api/user_web/forms/user"
	"mxshop_api/user_web/global"
	"mxshop_api/user_web/global/response"
	"mxshop_api/user_web/middlerwares"
	"mxshop_api/user_web/proto"
)

// GetUserList 获取用户列表
func GetUserList(c *gin.Context) {
	claims, _ := c.Get("claims")
	userInfo := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户:%d", userInfo.ID)

	pageNum := c.DefaultQuery("page_num", "1")
	pageNumInt, _ := strconv.Atoi(pageNum)
	pageSize := c.DefaultQuery("page_size", "10")
	pageSizeInt, _ := strconv.Atoi(pageSize)

	//限流
	e, b := sentinel.Entry("test", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		errReturn(c, http.StatusTooManyRequests, "请求过快")
		return
	}

	//请求服务
	rsp, err := global.UserSrvClient.GetUserList(context.WithValue(context.Background(), "ginContext", c), &proto.PageInfoRequest{
		PageNum:  uint32(pageNumInt),
		PageSize: uint32(pageSizeInt),
	})
	e.Exit()
	if err != nil {
		zap.S().Errorw("[GetUserList] 请求 [user_srv] 失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		userObj := response.UserResponse{
			Id:       value.Id,
			NickName: value.Nickname,
			Birthday: response.JsonTime(time.Unix(int64(value.Birthday), 0)),
			Gender:   strconv.FormatUint(uint64(value.Gender), 10),
			Mobile:   value.Mobile,
		}
		result = append(result, userObj)

		/*
			data := make(map[string]interface{})
			data["id"] = value.Id
			data["name"] = value.Nickname
			data["birthday"] = value.Birthday
			data["gender"] = value.Gender
			data["mobile"] = value.Mobile
			result = append(result,data)
		*/
	}
	c.JSON(http.StatusOK, result)
}

// PasswordLogin 用户密码登录
func PasswordLogin(c *gin.Context) {
	passwordLoginFrom := user.PasswordLoginFrom{}
	if err := bindJsonAndValidate(c, &passwordLoginFrom); !err {
		return
	}

	//if !store.Verify(passwordLoginFrom.CaptchaId,passwordLoginFrom.Captcha,true) {
	//	errReturn(c,http.StatusBadRequest,"验证码错误")
	//	return
	//}

	resp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: passwordLoginFrom.Mobile})
	if err != nil {
		zap.S().Errorw("[PasswordLogin] 请求 [user_srv.GetUserByMobile] 失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}
	checkResp, err := global.UserSrvClient.CheckPassword(context.Background(), &proto.CheckPasswordRequest{
		Md5Password: resp.Password,
		Password:    passwordLoginFrom.Password,
	})
	if err != nil {
		zap.S().Errorw("[PasswordLogin] 请求 [user_srv.CheckPassword] 失败", "msg", err.Error())
		errReturn(c, http.StatusBadRequest, "密码错误")
		return
	}
	if !checkResp.Success {
		errReturn(c, http.StatusBadRequest, "密码错误")
		return
	}
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(resp.Id),
		NickName:    resp.Nickname,
		AuthorityId: uint(resp.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + 60*60*24*30,
			Issuer:    "imooc",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		errReturn(c, http.StatusInternalServerError, "生成token失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         resp.Id,
		"nickName":   resp.Nickname,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}

func Register(c *gin.Context) {
	registerForm := user.RegisterFrom{}
	if ok := bindJsonAndValidate(c, &registerForm); !ok {
		return
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	defer rdb.Close()
	value, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil || value != registerForm.Code {
		errReturn(c, http.StatusBadRequest, "验证码错误")
		return
	}

	resp, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfoRequest{Mobile: registerForm.Mobile, Password: registerForm.Password, Nickname: registerForm.Mobile})
	if err != nil {
		zap.S().Errorw("[PasswordLogin] 请求 [user_srv.CreateUser] 失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(resp.Id),
		NickName:    resp.Nickname,
		AuthorityId: uint(resp.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + 60*60*24*30,
			Issuer:    "imooc",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		errReturn(c, http.StatusInternalServerError, "生成token失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         resp.Id,
		"nickName":   resp.Nickname,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}
