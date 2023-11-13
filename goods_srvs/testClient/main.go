package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"mxshop_srvs/goods_srvs/proto"
)

var (
	goodsClient proto.GoodsClient
	conn        *grpc.ClientConn
)

func Init() {
	var err error
	conn, err = grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		panic("conn error")
	}
	goodsClient = proto.NewGoodsClient(conn)

}

func main() {
	Init()
	brandsListRsp, err := goodsClient.BrandList(context.Background(), &proto.BrandFilterRequest{
		Pages: 2, PagePerNums: 5,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(brandsListRsp.Total)
	for _, brandsInfo := range brandsListRsp.Data {
		fmt.Println(brandsInfo.Name)
	}

	err = conn.Close()
	if err != nil {
		return
	}
}
