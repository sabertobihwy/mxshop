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

func TestCateBrand() {
	rsp, err := goodsClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{
		Pages: 2, PagePerNums: 5,
	})
	if err != nil {
		panic(err)
	}
	for _, v := range rsp.Data {
		fmt.Println(v.Brand)
	}
}
func TestGetCategoryBrandList() {
	rsp, err := goodsClient.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{
		Id: 130366,
	})
	if err != nil {
		panic(err)
	}
	for _, v := range rsp.Data {
		fmt.Println(v.Name)
	}
}
