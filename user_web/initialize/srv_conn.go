package initialize

import (
	"fmt"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"time"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"mxshop_api/user_web/global"
	"mxshop_api/user_web/proto"
	"mxshop_api/user_web/utils/otgrpc"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo

	//增加重试和超时
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithMax(3),
		grpc_retry.WithPerRetryTimeout(1 * time.Second),
		grpc_retry.WithCodes(codes.Unknown, codes.DeadlineExceeded, codes.Unavailable),
	}

	// 通过负载均衡器 去注册中心拿用户服务
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Errorw("[InitSrvConn] 连接 [user_srv] 失败", "msg", err.Error())
		return
	}
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}

func InitSrvConn2() {
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	userSrvHost := ""
	userSrvPort := 0

	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Errorf("连接注册中心失败:%s", err.Error())
		return
	}
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service==\"%s\"", global.ServerConfig.UserSrvInfo.Name))

	if err != nil {
		zap.S().Errorf("获取注册中心失败:%s", err.Error())
		return
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == "" {
		zap.S().Error("获取注册中心失败")
		return
	}

	fmt.Println(userSrvHost)

	// 拨号连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 [user_srv] 失败", "msg", err.Error())
		return
	}
	userSrvClient := proto.NewUserClient(userConn)
	//使用全局变量时存在的问题 1.当 用户服务下线,改端口,改IP时,会出现调用失败.  后续用负载均衡解决
	//使用全局变量,服务启动时就建立了连接,减少了tcp三次握手的消耗
	//一个连接 多个groutine使用,会存在竞争,  提高新能可能使用连接池  (https://github.com/processout/grpc-go-pool)
	global.UserSrvClient = userSrvClient
}
