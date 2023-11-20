package main

import (
	"context"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"

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
	_, err := stocksClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: 1,
		Num:     2,
	})
	if err != nil {
		panic(err)
	}
	_, err = stocksClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: 1,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Close()
	if err != nil {
		return
	}
}
