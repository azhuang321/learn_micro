package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"mytest/lock_test/proto"
	"sync"
)

var waitGroup sync.WaitGroup

func main() {
	conn, err := grpc.Dial("127.0.0.1:8000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewInventoryClient(conn)

	var goodsSlices [][]*proto.GoodsInvInfo
	for i := 0; i < 3; i++ {
		goodsSlices = append(goodsSlices, []*proto.GoodsInvInfo{
			&proto.GoodsInvInfo{GoodsId: 1, Num: 8},
			&proto.GoodsInvInfo{GoodsId: 2, Num: 8},
			&proto.GoodsInvInfo{GoodsId: 3, Num: 8},
		})
	}

	waitGroup.Add(len(goodsSlices))
	for _, val := range goodsSlices {
		go func() {
			fmt.Println(1)
			r, err := c.Sell(context.Background(), &proto.SellInfo{GoodsInfo: val})
			if err != nil {
				panic(err)
			}
			fmt.Println(r)
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
}
