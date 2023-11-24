package handler

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/goods_srvs/global"
	"mxshop_srvs/goods_srvs/model"
	"mxshop_srvs/goods_srvs/proto"
)

// 商品分类
func (s *GoodsServer) GetAllCategorysList(context.Context, *emptypb.Empty) (*proto.CategoryListResponse, error) {
	/*
		[
			{
				"id":xxx,
				"name":"",
				"level":1,
				"is_tab":false,
				"parent":13xxx,
				"sub_category":[
					"id":xxx,
					"name":"",
					"level":1,
					"is_tab":false,
					"sub_category":[]
				]
			}
		]
	*/
	var categorys []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	b, _ := json.Marshal(&categorys)
	return &proto.CategoryListResponse{JsonData: string(b)}, nil
}

// // 获取子分类
func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	var total int32
	var subs []model.Category
	var cate model.Category
	if result := global.DB.Where(&model.Category{BaseModel: model.BaseModel{ID: req.Id}}).First(&cate); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "not find this category")
	}
	var info = &proto.CategoryInfoResponse{
		Id:             cate.ID,
		Name:           cate.Name,
		ParentCategory: cate.ParentCategoryID,
		Level:          cate.Level,
		IsTab:          cate.IsTab,
	}

	result := global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Find(&subs)
	total = int32(result.RowsAffected)
	var subrsps []*proto.CategoryInfoResponse
	for _, value := range subs {
		subrsps = append(subrsps, &proto.CategoryInfoResponse{
			Id:             value.ID,
			Name:           value.Name,
			ParentCategory: value.ParentCategoryID,
			Level:          value.Level,
			IsTab:          value.IsTab,
		})
	}
	return &proto.SubCategoryListResponse{
		Total:        total,
		Info:         info,
		SubCategorys: subrsps,
	}, nil
}

//CreateCategory(ctx context.Context, in *CategoryInfoRequest, opts ...grpc.CallOption) (*CategoryInfoResponse, error)
//DeleteCategory(ctx context.Context, in *DeleteCategoryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
//UpdateCategory(ctx context.Context, in *CategoryInfoRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
//
