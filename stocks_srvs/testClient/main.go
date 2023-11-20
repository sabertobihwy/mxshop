package main

import (
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"sync"

	"mxshop_srvs/stocks_srvs/proto"
)

var (
	stocksClient proto.StocksClient
	conn         *grpc.ClientConn
)

func Init() {
	var err error
	conn, err = grpc.Dial("consul://192.168.2.112:8500/stocks_srvs?wait=14s&tag=mxshop",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		panic("conn error")
	}
	stocksClient = proto.NewStocksClient(conn)

}

func main() {
	Init()

	//TestInvDetail()
	//TestSetInv()
	//TestSell()
	//TestReBack()
	var i int32
	var wg sync.WaitGroup
	wg.Add(5)
	defer conn.Close()
	defer wg.Wait()
	//defer conn.Close()
	for i = 0; i < 5; i++ {
		go TestSell(&wg)
	}
}

func TestSetInv(i int32, n int32) {
	_, err := stocksClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: i,
		Num:     n,
	})
	if err != nil {
		panic(err)
	}
}

func TestInvDetail() {
	num, err := stocksClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: 1,
	})
	fmt.Println(num)
	if err != nil {
		panic(err)
	}
}

func TestSell(wg *sync.WaitGroup) {
	defer wg.Done()
	num, err := stocksClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInvInfo: []*proto.GoodsInvInfo{
			{
				GoodsId: 421,
				Num:     10,
			},
			{
				GoodsId: 422,
				Num:     10,
			},
			{
				GoodsId: 423,
				Num:     10,
			},
		},
	})
	fmt.Println(num)
	if err != nil {
		panic(err)
	}
}
func TestReBack() {
	num, err := stocksClient.Reback(context.Background(), &proto.SellInfo{
		GoodsInvInfo: []*proto.GoodsInvInfo{
			{
				GoodsId: 1,
				Num:     10,
			},
			{
				GoodsId: 2,
				Num:     10,
			},
		},
	})
	fmt.Println(num)
	if err != nil {
		panic(err)
	}
}
