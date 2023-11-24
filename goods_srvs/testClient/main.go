package main

import (
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"mxshop_srvs/goods_srvs/proto"
)

var (
	goodsClient proto.GoodsClient
	conn        *grpc.ClientConn
)

func Init() {
	var err error
	conn, err = grpc.Dial("consul://192.168.2.112:8500/goods_srvs?wait=14s&tag=mxshop",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		panic("conn error")
	}
	goodsClient = proto.NewGoodsClient(conn)

}

func main() {
	Init()
	TestSubCategory()
	//TestBrand()
	conn.Close()
}
