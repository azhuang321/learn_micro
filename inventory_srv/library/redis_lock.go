package library

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisLock struct {
	Key  string
	UUid string
}

var ctx = context.Background()

func (l *RedisLock) Acquire() bool {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()
	val, err := rdb.SetNX(ctx, l.Key, l.UUid, 15).Result() //这样能保证原子操作 否则先判断再获取 存在两个步骤不能保证原子操作
	if err != nil {
		fmt.Println("加锁失败", l.Key)
		return false
	} else if val {
		// 这里防止脚本执行 超过当前设置的过期时间,所以需要启动一个协程 当做看门狗 来续租这个过期时间
		// https://github.com/ionelmc/python-redis-lock/blob/master/src/redis_lock/__init__.py#L178  第三方库实现的redis_lock
		// 第三方库 实现更完善 https://github.com/go-locks/distlock (中文版)
		//https://github.com/go-redsync/redsync (英文版)

		fmt.Println("加锁成功", l.Key)
		return true
	} else {
		// 没有获取到锁 应该阻塞
		for {
			time.Sleep(time.Second)
			val, err := rdb.SetNX(ctx, l.Key, l.UUid, 15).Result()
			if val && err == nil {
				fmt.Println("等待加锁成功", l.Key)
				return true
			}
			fmt.Println("获取锁失败 等待获取", l.Key)
		}
	}
}

// 这里有个问题 因为是先获取值 再去设置值 导致这不是一个原子操作 所以 为了解决这个问题 可以使用lua 脚本 使其成为一个原子操作
func (l *RedisLock) Release() bool {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	val, err := rdb.Get(ctx, l.Key).Result()
	if err != nil && err != redis.Nil {
		fmt.Println("解锁失败", l.Key)
		return false
	}
	if val != l.UUid {
		fmt.Println("不是解锁属于自己的锁", l.Key)
		return false
	}

	err = rdb.Del(ctx, l.Key).Err()
	if err != nil {
		fmt.Println("解锁失败", l.Key)
		return false
	}
	fmt.Println("解锁成功", l.Key)
	return true
}
