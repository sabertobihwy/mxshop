package main

import (
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"mxshop_srvs/order_srvs/proto"
)

var (
	orderClient proto.OrderClient
	conn        *grpc.ClientConn
)

func Init() {
	var err error
	conn, err = grpc.Dial("consul://192.168.2.112:8500/order_srvs?wait=14s&tag=mxshop",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		panic("conn error")
	}
	orderClient = proto.NewOrderClient(conn)

}

func main() {
	Init()
	TestCreateCart(1, 1, 1)
}

func TestCreateCart(userId int32, nums int32, goodsId int32) {
	rsp, err := orderClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  userId,
		Nums:    nums,
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(rsp.Id)
}
