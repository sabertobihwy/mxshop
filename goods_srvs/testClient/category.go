package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"mxshop_srvs/goods_srvs/proto"
)

func TestCategory() {
	result, err := goodsClient.GetAllCategorysList(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println(result.JsonData)

}

func TestSubCategory() {
	result, err := goodsClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: 130364,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(result.SubCategorys)

}
