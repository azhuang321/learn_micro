package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"user_srv/database"
	"user_srv/initialize"
	"user_srv/model"
	"user_srv/proto/gen/coupons_pb"
	"user_srv/proto/gen/user_pb"
	"user_srv/utils"
)

type UserService struct{}

func convertModelUserToResponseUser(user model.User) *user_pb.UserInfoResponse {
	respUser := user_pb.UserInfoResponse{}
	respUser.Id = int32(user.ID)
	respUser.Password = user.Password
	respUser.Mobile = user.Mobile
	respUser.Nickname = user.Nickname
	return &respUser
}

func (u *UserService) GetUserList(ctx context.Context, request *user_pb.PageInfoRequest) (resp *user_pb.UserListResponse, err error) {
	var users []model.User
	db, err := database.GetDB()
	if err != nil {
		zap.S().Errorf("数据库错误:%s\n", err.Error())
		return nil, status.Error(codes.DataLoss, "获取数据库出错")
	}
	var count int64
	db.Model(&users).Count(&count)
	resp = &user_pb.UserListResponse{}
	resp.Total = int32(count)

	var page uint32 = 1
	var pageNum uint32 = 10
	if request.PageSize > 0 {
		pageNum = request.PageSize
	}
	if request.PageNum > 1 {
		page = request.PageNum
	}
	offset := pageNum * (page - 1)
	db.Offset(int(offset)).Limit(int(pageNum)).Find(&users)
	for _, value := range users {
		userInfoResp := convertModelUserToResponseUser(value)
		resp.Data = append(resp.Data, userInfoResp)
	}

	return resp, nil
}

type DemoListener struct {
	localTrans       *sync.Map
	transactionIndex int32
}

func NewDemoListener() *DemoListener {
	return &DemoListener{
		localTrans: new(sync.Map),
	}
}

var localExecDict = make(map[string]map[string]interface{})

// 本地事务 基于回调,发送消息成功后 立即调用此回调方法,依据返回的状态 来确认 发送的消息是否可进行消费
func (dl *DemoListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	msgBody := msg.Body
	msgs := make(map[string]interface{})
	_ = json.Unmarshal(msgBody, &msgs)
	mobile := msgs["mobile"].(string)
	couponsId, _ := strconv.Atoi(msgs["couponsId"].(string))
	num, _ := strconv.Atoi(msgs["num"].(string))

	var user model.User
	db, err := database.GetDB()
	if err != nil {
		zap.S().Errorf("数据库错误:%s\n", err.Error())
		localExecDict[mobile]["code"] = codes.DataLoss
		localExecDict[mobile]["err"] = status.Error(codes.DataLoss, "获取数据库出错")
		return primitive.RollbackMessageState
	}

	db.Where("mobile = ?", msgs["mobile"]).First(&user)
	if user.ID > 0 {
		localExecDict[mobile]["code"] = codes.AlreadyExists
		localExecDict[mobile]["err"] = status.Error(codes.AlreadyExists, "用户已存在")
		return primitive.RollbackMessageState
	}

	user = model.User{Nickname: msgs["nickname"].(string), Password: utils.MD5(msgs["password"].(string)), Mobile: msgs["mobile"].(string)}
	result := db.Create(&user)
	if result.Error != nil {
		localExecDict[mobile]["code"] = codes.Unknown
		localExecDict[mobile]["err"] = status.Error(codes.Unknown, "创建用户失败,未知原因")
		return primitive.RollbackMessageState
	}
	//注册送积分
	_, err = initialize.CouponsClient.SendCouponsToUser(context.Background(), &coupons_pb.SendCouponsToUserRequest{
		Mobile:    mobile,
		CouponsId: uint32(couponsId),
		Num:       uint32(num),
	})

	if err != nil {
		localExecDict[mobile]["code"] = codes.Internal
		localExecDict[mobile]["err"] = status.Error(codes.Internal, "赠送积分失败")
		/*
			调用失败问题比较复杂
		*/
		grpcErr, _ := status.FromError(err)
		if grpcErr.Code() == codes.Unknown || grpcErr.Code() == codes.DeadlineExceeded {
			//返回未知错误和 超时错误 就需要消息回查
			return primitive.CommitMessageState
		} else {
			return primitive.RollbackMessageState
		}
	}

	userInfoResp := convertModelUserToResponseUser(user)
	localExecDict[mobile]["code"] = codes.OK
	localExecDict[mobile]["err"] = nil
	localExecDict[mobile]["resp"] = userInfoResp

	//发送延时消息 取消优惠券
	rlog.SetLogLevel("error")
	p, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"172.18.0.1:9876"})),
		producer.WithRetry(2),
	)
	err = p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	//生成每笔赠送积分单号

	delayMsg := map[string]string{
		"couponsId": "1",
		"mobile":    mobile,
		"num":       msgs["num"].(string),
	}
	delayMsgBody, _ := json.Marshal(delayMsg)

	sendMsg := primitive.NewMessage("test1", delayMsgBody)
	sendMsg.WithDelayTimeLevel(3)
	res, err := p.SendSync(context.Background(), sendMsg)

	if err != nil {
		return primitive.CommitMessageState
	}

	if res.Status != primitive.SendOK {
		return primitive.CommitMessageState
	}

	fmt.Printf("send status:%+v,message_id:%v", res.Status, res.MsgID)

	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}

	return primitive.RollbackMessageState
}

// 当长时间没有得到此事务消息的本地事务状态,就会进行消息回查,当出现宕机的异常情况 消息没有反馈明确状态时,通过回查来确认用户是否创建成功
func (dl *DemoListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	msgBody := msg.Body
	msgs := make(map[string]interface{})
	_ = json.Unmarshal(msgBody, &msgs)
	//pointsSn := msgs["pointsSn"].(string)

	//查询本地数据库,看一下points_sn是否赠送成功
	//todo 设计逻辑 调用优惠券服务查询是否赠送成功
	if true {
		return primitive.RollbackMessageState //不归还
	} else {
		return primitive.CommitMessageState //归还
	}
}

/*CreateUser

调用积分归还接口的问题：
	1. 在调用接口之前，程序出现了异常
	2. 在调用之前非程序异常，都可能会导致接口调用失败
	3. 调用接口的时候网络出现了抖动 - 幂等性机制

注册活动
注册送积分 1  总积分100 送完即止

场景: 当用户注册失败时,需要归还当前注册用户前冻结的积分到总积分100积分中去
分布式赠送积分逻辑解决方案: (利用rocketmq的可靠事务消息)
	1.在注册用户赠送积分前,准备好一个 rocketmq 的half消息
	2.注册用户
		*当注册用户失败:确认消息 - 积分服务 可以确定归还积分
		*当用户注册成功:
			- 请求积分服务出现异常
			- 积分服务返回出现异常: 假如取消 消息时,积分服务收到取消积分,但是在返回结果时出现了意想不到的异常(如:网络抖动,宕机等),导致超时机制认为这个调用失败
				则确认消息 - 确认消息 确认消息成功 则积分服务可以退还刚才的积分(但是积分服务还是需要记录一下之前的积分是否已经归还)
			- 本地事务提交执行异常:
				还是确认消息
			- 本地宕机:
				rocketmq就会启动回查机制 -
					本地的回查业务： 通过消息中的订单号在本地查询数据库中是否有数据
						取消消息
					查询到没有数据：
						确认消息

 # 创建用户 赠送积分
	# 原本比较简单的逻辑应该是： 1. 本地开始half消息 - 赠送积分 2. 执行本地事务 3. 确定应该确认消息还是回滚消息
	# 1. 基于可靠消息的最终一致性 只确保自己发送出去的消息是可靠的， 不能确保消费者能正确的执行
	# 2. 积分服务 - 这里有一个隐含的点： 你的消费者必要保证能成功
	# 3. 有限积分活动比较特殊 - 积分是有限的 如果本地事务执行失败应该调用归还积分 - 不使用TCC ：1. 并发没有那么高 2. 很复杂
*/
func (u *UserService) CreateUser(ctx context.Context, info *user_pb.CreateUserInfoRequest) (*user_pb.UserInfoResponse, error) {
	localExecDict[info.Mobile] = map[string]interface{}{}
	rlog.SetLogLevel("error")
	p, _ := rocketmq.NewTransactionProducer(
		NewDemoListener(),
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"172.18.0.1:9876"})),
		producer.WithRetry(1),
	)
	err := p.Start()

	if err != nil {
		fmt.Printf("start producer error: %s\n", err.Error())
		os.Exit(1)
	}
	//生成每笔赠送积分单号

	msg := map[string]string{
		"couponsId": "1",
		"mobile":    info.Mobile,
		"nickname":  info.Nickname,
		"password":  info.Password,
		"num":       info.CouponsNum,
	}
	msgBody, _ := json.Marshal(msg)

	res, err := p.SendMessageInTransaction(context.Background(), primitive.NewMessage("test", msgBody))

	if err != nil {
		fmt.Printf("send message error: %s\n", err)
	}
	if res.Status != primitive.SendOK {
		return nil, status.Error(codes.Internal, "赠送积分失败")
	}

	fmt.Printf("send status:%+v,message_id:%v", res.Status, res.MsgID)

	//todo 改为chan
	for {
		if _, ok := localExecDict[info.Mobile]; ok {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
	//todo 删除全局表里中的
	if localExecDict[info.Mobile]["code"].(codes.Code) == codes.OK {
		return localExecDict[info.Mobile]["resp"].(*user_pb.UserInfoResponse), nil
	}
	return nil, localExecDict[info.Mobile]["err"].(error)
}
