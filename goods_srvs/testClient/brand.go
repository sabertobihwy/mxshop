package main

import (
	"context"
	"fmt"
	"mxshop_srvs/goods_srvs/proto"
)

func TestBrand() {
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
}
