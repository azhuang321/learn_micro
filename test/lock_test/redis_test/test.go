package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisLock struct {
	Key string
}

var ctx = context.Background()

func (l *redisLock) acquire() bool {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()
	val, err := rdb.SetNX(ctx, l.Key, 1, 0).Result() //这样能保证原子操作 否则先判断再获取 存在两个步骤不能保证原子操作
	if err != nil {
		return false
	} else if val {
		return true
	} else {
		// 没有获取到锁 应该阻塞
		for {
			time.Sleep(time.Second)
			val, err := rdb.SetNX(ctx, l.Key, 1, 0).Result()
			if val && err == nil {
				return true
			}
			fmt.Println("获取锁失败 等待获取")
		}
	}
}

func (l *redisLock) Release() bool {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	err := rdb.Del(ctx, l.Key).Err()
	if err != nil {
		return false
	}
	return true
}

func main() {
	goodsId := 1
	redisLock := redisLock{Key: fmt.Sprintf("lock_goods_%d", goodsId)}
	redisLock.acquire()

}
