package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"gorm.io/gorm"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"coupons_srv/database"
	"coupons_srv/model"
	"coupons_srv/proto/gen/coupons_pb"
)

type CouponsService struct{}

var i int = 0

func (p CouponsService) SendCouponsToUser(ctx context.Context, request *coupons_pb.SendCouponsToUserRequest) (*emptypb.Empty, error) {
	i++
	db, err := database.GetDB()
	if err != nil {
		zap.S().Errorf("数据库错误:%s\n", err.Error())
		return nil, status.Error(codes.DataLoss, "获取数据库出错")
	}
	//先查询是否有对应优惠券
	// todo 分布式锁
	var coupons model.Coupons
	if err := db.First(&coupons, request.CouponsId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "没有找到优惠券")
		}
		return nil, status.Error(codes.Internal, "获取数据库出错")
	}
	if coupons.Num < request.Num {
		return nil, status.Error(codes.OutOfRange, "优惠券不足")
	}

	tx := db.Begin()
	//添加扣取记录,以便在归还的时候查询
	var couponsHistory model.CouponsHistory
	couponsHistory.Num = request.Num
	couponsHistory.Mobile = request.Mobile
	couponsHistory.Status = 1
	if err := tx.Create(&couponsHistory).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, "获取数据库出错")
	}
	//扣取
	if err := tx.Model(&coupons).Update("num", gorm.Expr("num - ?", request.Num)).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, "获取数据库出错")
	}

	//插入赠送记录
	var couponsUser model.CouponsUser
	couponsUser.Mobile = request.Mobile
	couponsUser.Num = request.Num
	if err = tx.Create(&couponsUser).Error; err != nil {
		tx.Rollback()
		zap.S().Errorf("数据库错误:%s\n", err.Error())
		return nil, status.Error(codes.Internal, "插入失败,未知错误")
	}
	tx.Commit()

	if i <= 1 {
		fmt.Println("正常逻辑")
		return &emptypb.Empty{}, nil

	}
	fmt.Println("测试返回异常逻辑")
	return nil, status.Error(codes.Unknown, "获取数据库出错")
}

// DealMsg 扣减库存
func DealMsg() func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	return func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			msgBody := msgs[i].Body
			msg := make(map[string]interface{})
			_ = json.Unmarshal(msgBody, &msg)
			fmt.Println(msg)
			couponsId := msg["couponsId"].(string)
			mobile := msg["mobile"].(string)
			num := msg["num"].(string)
			db, err := database.GetDB()
			if err != nil {
				zap.S().Errorf("数据库错误:%s\n", err.Error())
				return consumer.ConsumeSuccess, nil
			}

			//为了防止没有赠送反而返还了积分的情况,这里我们要先查询有没有积分扣减记录
			var couponsHistory model.CouponsHistory
			err = db.Where("mobile = ?", mobile).Where("status = ?", 1).First(&couponsHistory).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return consumer.ConsumeSuccess, nil
				}
				return consumer.ConsumeRetryLater, nil
			}
			tx := db.Begin()
			var coupons model.Coupons
			if err := tx.Model(&coupons).Where("id = ?", couponsId).Update("num", gorm.Expr("num + ?", num)).Error; err != nil {
				fmt.Println(err)
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
			if err := tx.Model(&couponsHistory).Update("status", 2).Error; err != nil {
				fmt.Println(err)
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
			var couponsUser model.CouponsUser
			couponsUser.Mobile = mobile
			if err := tx.Where("mobile = ?", mobile).Delete(&couponsUser).Error; err != nil {
				fmt.Println(err)
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}

			tx.Commit()
			return consumer.ConsumeSuccess, nil
		}
		return consumer.ConsumeSuccess, nil
	}
}
