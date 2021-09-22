package handler

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"inventory_srv/database"
	"inventory_srv/library"
	"inventory_srv/model"
	"inventory_srv/proto"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type InventoryService struct{}

func (i InventoryService) SetInv(ctx context.Context, info *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	db, err := database.GetDB()
	if err != nil {
		zap.S().Errorf("数据库错误:%s\n", err.Error())
		return nil, status.Error(codes.DataLoss, "获取数据库出错")
	}
	var inv model.Inventory
	db.Where("goods = ?", info.GoodsId).First(&inv)
	if inv.ID > 0 {
		inv.Stocks = int(info.Num)
		db.Save(&inv)
	} else if info.GoodsId > 0 {
		inv.Goods = int(info.GoodsId)
		inv.Stocks = int(info.Num)
		db.Create(&inv)
	}
	return &emptypb.Empty{}, nil
}

func (i InventoryService) InvDetail(ctx context.Context, info *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	db, err := database.GetDB()
	if err != nil {
		zap.S().Errorf("数据库错误:%s\n", err.Error())
		return nil, status.Error(codes.DataLoss, "获取数据库出错")
	}
	var inv model.Inventory
	db.Where("goods = ?", info.GoodsId).First(&inv)
	if inv.ID <= 0 {
		return nil, status.Error(codes.NotFound, "暂无数据")
	}
	return &proto.GoodsInvInfo{GoodsId: int32(inv.Goods), Num: int32(inv.Stocks)}, nil
}

// todo redis 锁 这里有bug  不同事务到达的时间不一致 在上一次事务未提交前 下一次事务又到达了 那就不会读取到上一次修改但未提交的数据  导致超卖
func (i InventoryService) Sell(ctx context.Context, info *proto.SellInfo) (*emptypb.Empty, error) {
	db, err := database.GetDB()
	if err != nil {
		zap.S().Errorf("数据库错误:%s\n", err.Error())
		return nil, status.Error(codes.DataLoss, "获取数据库出错")
	}
	tx := db.Begin()
	for _, val := range info.GoodsInfo {
		u2 := uuid.NewV4()
		redisLock := library.RedisLock{Key: "local_goods_" + strconv.Itoa(int(val.GoodsId)), UUid: fmt.Sprintf("%s", u2)}
		redisLock.Acquire()

		var goodsInv model.Inventory
		tx.Where("goods = ?", val.GoodsId).First(&goodsInv)
		fmt.Printf("%+v\n", goodsInv)
		//模拟业务耗时
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

		if goodsInv.ID <= 0 || goodsInv.Stocks < int(val.Num) {
			redisLock.Release()
			tx.Rollback()
			return nil, status.Error(codes.ResourceExhausted, "库存不足")
		}

		//todo 可能会引起数据不一致
		if err := tx.Model(&goodsInv).Update("stocks", gorm.Expr("stocks - ? ", val.Num)).Error; err != nil {
			redisLock.Release()
			tx.Rollback()
			return nil, status.Error(codes.ResourceExhausted, "库存不足")
		}
		redisLock.Release()
	}
	tx.Commit()
	//加锁分开加 ,释放锁 需要在commit之后释放
	return &emptypb.Empty{}, nil
}

// 基于 mysql的version 实现乐观锁
func (i InventoryService) Sell2(ctx context.Context, info *proto.SellInfo) (*emptypb.Empty, error) {
	db, err := database.GetDB()
	if err != nil {
		zap.S().Errorf("数据库错误:%s\n", err.Error())
		return nil, status.Error(codes.DataLoss, "获取数据库出错")
	}
	tx := db.Begin()
	for _, val := range info.GoodsInfo {
		for {
			var goodsInv model.Inventory
			db.Where("goods = ?", val.GoodsId).First(&goodsInv)
			fmt.Println("当前版本号:", goodsInv.Version)
			//模拟业务耗时
			rand.Seed(time.Now().UnixNano())
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

			if goodsInv.ID <= 0 || goodsInv.Stocks < int(val.Num) {
				tx.Rollback()
				return nil, status.Error(codes.ResourceExhausted, "库存不足")
			}

			if affectedNum := tx.Model(&goodsInv).Where("version = ?", goodsInv.Version).Updates(map[string]interface{}{"stocks": gorm.Expr("stocks - ? ", val.Num), "version": gorm.Expr("version + ? ", 1)}).RowsAffected; affectedNum == 0 {
				fmt.Println("1111111111111111111111111111111111111111111111")
			} else {
				break
			}
		}
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

var lock sync.Mutex

// 悲观锁 直接每次加锁 (效率低)
func (i InventoryService) Sell1(ctx context.Context, info *proto.SellInfo) (*emptypb.Empty, error) {
	db, err := database.GetDB()
	if err != nil {
		zap.S().Errorf("数据库错误:%s\n", err.Error())
		return nil, status.Error(codes.DataLoss, "获取数据库出错")
	}
	tx := db.Begin()
	lock.Lock()
	for _, val := range info.GoodsInfo {
		var goodsInv model.Inventory
		db.Where("goods = ?", val.GoodsId).First(&goodsInv)
		//模拟业务耗时
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

		if goodsInv.ID <= 0 || goodsInv.Stocks < int(val.Num) {
			tx.Rollback()
			return nil, status.Error(codes.ResourceExhausted, "库存不足")
		}
		//todo 可能会引起数据不一致
		if err := tx.Model(&goodsInv).Update("stocks", gorm.Expr("stocks - ? ", val.Num)).Error; err != nil {
			tx.Rollback()
			return nil, status.Error(codes.ResourceExhausted, "库存不足")
		}
	}
	lock.Unlock()
	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (i InventoryService) Reback(ctx context.Context, info *proto.SellInfo) (*emptypb.Empty, error) {
	db, err := database.GetDB()
	if err != nil {
		zap.S().Errorf("数据库错误:%s\n", err.Error())
		return nil, status.Error(codes.DataLoss, "获取数据库出错")
	}
	tx := db.Begin()
	for _, val := range info.GoodsInfo {
		var goodsInv model.Inventory
		db.Where("goods = ?", val.GoodsId).First(&goodsInv)
		if goodsInv.ID <= 0 {
			tx.Rollback()
			return nil, status.Error(codes.NotFound, "商品不存在")
		}
		//todo 可能会引起数据不一致
		goodsInv.Stocks += int(val.Num)
		if err := tx.Save(&goodsInv).Error; err != nil {
			tx.Rollback()
			return nil, status.Error(codes.ResourceExhausted, "库存不足")
		}
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
