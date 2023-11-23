package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"mxshop_srvs/order_srvs/global"
	"mxshop_srvs/order_srvs/model"
	"mxshop_srvs/order_srvs/proto"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

func Paginate(pg, pgSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch {
		case pgSize > 100:
			pgSize = 100
		case pgSize <= 0:
			pgSize = 10
		}
		offset := (pg - 1) * pgSize
		return db.Offset(offset).Limit(pgSize)
	}
}

func (*OrderServer) CartItemList(c context.Context, u *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var rsp proto.CartItemListResponse
	var shopCarts []model.ShoppingCart
	if result := global.DB.Where(&model.ShoppingCart{User: u.Id}).Find(&shopCarts); result.Error != nil {
		return nil, result.Error
	} else {
		rsp.Total = int32(result.RowsAffected)
	}
	for _, shopCart := range shopCarts {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id:      shopCart.ID,
			UserId:  shopCart.User,
			GoodsId: shopCart.Goods,
			Nums:    shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}
	return &rsp, nil
}
func (*OrderServer) CreateCartItem(ctx context.Context, c *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	//  goods already in the cart: update; not in: create
	var cart model.ShoppingCart
	if result := global.DB.Where(&model.ShoppingCart{Goods: c.GoodsId, User: c.UserId}).First(&cart); result.RowsAffected == 1 {
		// update
		cart.Nums += c.Nums
	} else {
		cart.Goods = c.GoodsId
		cart.User = c.UserId
		cart.Checked = false
		cart.Nums = c.Nums
	}
	global.DB.Save(&cart)
	return &proto.ShopCartInfoResponse{Id: cart.ID}, nil
}
func (*OrderServer) UpdateCartItem(ctx context.Context, cart *proto.CartItemRequest) (*emptypb.Empty, error) {
	// update check state and number
	var shopCart model.ShoppingCart
	if result := global.DB.First(&shopCart, cart.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "Record not exists")
	}
	if cart.Nums > 0 {
		shopCart.Nums = cart.Nums
	}
	shopCart.Checked = cart.Checked
	global.DB.Save(&shopCart)
	return &emptypb.Empty{}, nil
}
func (*OrderServer) DeleteCartItem(ctx context.Context, cart *proto.CartItemRequest) (*emptypb.Empty, error) {
	var shopCart model.ShoppingCart
	if result := global.DB.Delete(&shopCart, cart.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "Record not exists")
	}
	return &emptypb.Empty{}, nil
}
func (*OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var orders []model.OrderInfo
	var total int64
	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	var rsp proto.OrderListResponse
	rsp.Total = int32(total)
	// pagination query
	if result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&orders); result != nil {
		return nil, status.Errorf(codes.Internal, "query for orderlist error")
	}
	for _, order := range orders {
		rsp.Data = append(rsp.Data, &proto.OrderInfoResponse{
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
		})
	}
	return &rsp, nil
}
func (*OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	// for backend management system, req needs orderid; for web, req needs orderid and userid
	// GORM will ignore zero condition
	var order model.OrderInfo
	result := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).Find(&order)
	if result != nil || result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "order not exists")
	}
	var rsp proto.OrderInfoDetailResponse
	rsp.OrderInfo.Id = order.ID
	rsp.OrderInfo.OrderSn = order.OrderSn
	rsp.OrderInfo.UserId = order.User
	rsp.OrderInfo.Address = order.Address
	rsp.OrderInfo.Post = order.Post
	rsp.OrderInfo.Status = order.Status
	rsp.OrderInfo.PayType = order.PayType
	rsp.OrderInfo.Total = order.OrderMount
	rsp.OrderInfo.Name = order.SignerName
	rsp.OrderInfo.Mobile = order.SingerMobile
	var goods []model.OrderGoods
	if result = global.DB.Where(&model.OrderGoods{Order: order.ID}).Find(&goods); result != nil {
		return nil, result.Error
	}
	for _, good := range goods {
		rsp.Goods = append(rsp.Goods, &proto.OrderItemResponse{
			OrderId:    good.Order,
			GoodsId:    good.Goods,
			GoodsName:  good.GoodsName,
			GoodsImage: good.GoodsImage,
			GoodsPrice: good.GoodsPrice,
			Nums:       good.Nums,
		})
	}

	return &rsp, nil
}
