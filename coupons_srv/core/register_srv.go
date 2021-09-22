package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"net"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	. "coupons_srv/config"
	"coupons_srv/handler"
	"coupons_srv/proto/gen/coupons_pb"
	"coupons_srv/register"
	"coupons_srv/utils"
)

func RegisterService() {
	var err error
	register.SrvRegister, err = register.NewConsulRegister()
	if err != nil {
		zap.S().Errorf("注册服务中心失败:%s", err.Error())
	}

	gs := grpc.NewServer()
	//注册用户服务
	couponsSrv := &handler.CouponsService{}
	coupons_pb.RegisterCouponsServer(gs, couponsSrv)

	//注册健康检查服务
	healthCheckSrv := &handler.HealthCheckSrv{Status: grpc_health_v1.HealthCheckResponse_SERVING, Reason: "running"}
	grpc_health_v1.RegisterHealthServer(gs, healthCheckSrv)

	args := utils.GetArgs()
	//lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", args["host"], args["port"]))
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 9000))

	if err != nil {
		zap.S().Fatalf("监听端口出现错误:%s", err.Error())
		return
	}

	go func() {
		zap.S().Infof("启动服务成功:%s:%d", args["host"], args["port"])
		if err := gs.Serve(lis); err != nil {
			zap.S().Errorf("启动服务失败:%s\n", err.Error())
		}
	}()

	zap.S().Info("开始注册服务中心....")
	if register.SrvRegister.Register(Config.ProjectName, Config.Uuid, args["host"].(string), args["port"].(int), Config.Consul.Tags, nil) {
		zap.S().Info("注册服务中心成功")
	}

	rlog.SetLogLevel("error")
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName("testGroup"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"172.18.0.1:9876"})),
	)

	err = c.Subscribe("test", consumer.MessageSelector{}, handler.DealMsg())
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	nc, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName("testGroup1"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"172.18.0.1:9876"})),
	)

	err = nc.Subscribe("test1", consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i, _ := range msgs {
			msgBody := msgs[i].Body
			msg := make(map[string]interface{})
			_ = json.Unmarshal(msgBody, &msg)
			fmt.Println(msg)
			//todo  消费消息  定时取消未领取优惠券

		}

		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = nc.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

}
