package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"mxshop_srvs/order_srvs/global"
	"mxshop_srvs/order_srvs/model"
	"mxshop_srvs/order_srvs/proto"
	"os"
	"time"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

// y m d h m s + user_id+ random 2 digits
func GenerateOrderSn(userid int32) string {
	now := time.Now()
	rand.Seed(uint64(now.UnixNano()))
	return fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userid, rand.Intn(90)+10,
	)
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
	var odrgoods []model.OrderGoods
	if result = global.DB.Where(&model.OrderGoods{Order: order.ID}).Find(&odrgoods); result != nil {
		return nil, result.Error
	}
	for _, good := range odrgoods {
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

type OrderListener struct {
	Code       codes.Code
	Detail     string
	Id         int32
	OrderMount float32
}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	fmt.Printf("执行本地逻辑\n")
	order := model.OrderInfo{}
	_ = json.Unmarshal(msg.Body, &order)

	// 1. get good_ids
	var carts []model.ShoppingCart
	var cartMap = make(map[int32]int32)
	if result := global.DB.Where(&model.ShoppingCart{User: order.User, Checked: true}).Find(&carts); result.RowsAffected == 0 {
		o.Code = codes.InvalidArgument
		o.Detail = "no checked items!"
		return primitive.RollbackMessageState
	}
	var good_ids []int32
	for _, cart := range carts {
		good_ids = append(good_ids, cart.Goods)
		cartMap[cart.Goods] = cart.Nums
	}
	// 2. Calling microservices across servers: good_ids -> goodsInfo -> generate ordergoods
	GoodsListResponse, err := global.GoodsClient.BatchGetGoods(context.Background(),
		&proto.BatchGoodsIdInfo{
			Id: good_ids,
		})
	if err != nil {
		o.Code = codes.Internal
		o.Detail = "invoking Goods-Service fails"
		return primitive.RollbackMessageState
	}
	// insert in batch
	var orderGoods []*model.OrderGoods
	var price float32
	var goodsInvInfo []*proto.GoodsInvInfo
	for _, goods := range GoodsListResponse.Data {
		price += goods.ShopPrice * float32(cartMap[goods.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      goods.Id,
			GoodsName:  goods.Name,
			GoodsImage: goods.GoodsFrontImage,
			GoodsPrice: goods.ShopPrice,
			Nums:       cartMap[goods.Id],
		})
		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: goods.Id,
			Num:     cartMap[goods.Id],
		})
	}

	// 3. reduce inventory
	if _, err = global.StocksClient.Sell(context.Background(), &proto.SellInfo{GoodsInvInfo: goodsInvInfo}); err != nil {
		// identify if it is network problem, if code is not xxx -> internet problem
		// ???
		o.Code = codes.ResourceExhausted
		o.Detail = "reducing stocks fail"
		return primitive.RollbackMessageState
	}
	// 4. order
	tx := global.DB.Begin()
	order.OrderMount = price
	if result := tx.Save(&order); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "invoking Save fails"
		return primitive.CommitMessageState

	}
	o.OrderMount = price
	o.Id = order.ID
	for _, og := range orderGoods {
		og.Order = order.ID
	}
	// insert ordergoods in batch
	if result := tx.CreateInBatches(&orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "invoking CreateInBatches fails"
		return primitive.CommitMessageState
	}

	// 5. del from cart
	if result := tx.Where(&model.ShoppingCart{User: order.User, Checked: true}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "delete from cart fails"
		return primitive.CommitMessageState
	}

	tx.Commit()
	fmt.Printf("执行完毕\n")
	return primitive.RollbackMessageState
}

func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Printf("消息回查\n")
	order := model.OrderInfo{}
	_ = json.Unmarshal(msg.Body, &order)
	if result := global.DB.Where(&model.OrderInfo{OrderSn: order.OrderSn}).First(&order); result.RowsAffected == 0 {
		// 不清楚库存是否扣减，但选择commit让消费端执行reback() -> 要求消费端执行reback()的时候进行幂等判断
		return primitive.CommitMessageState
	}
	return primitive.RollbackMessageState // order成功说明库存必然成功
}

/*
4. CreateOrder:

 1. get checked goods from cart

 2. add price - visit GoodService

 3. reduce stocks - visit StockService

 4. order basic info

 5. del from cart
*/
func (*OrderServer) CreateOrder(c context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	ol := OrderListener{}
	p, _ := rocketmq.NewTransactionProducer(
		&ol,
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.2.112:9876"})),
		producer.WithRetry(2),
	)
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	order := model.OrderInfo{
		User:         req.UserId,
		OrderSn:      GenerateOrderSn(req.UserId),
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
	}
	jsonStr, _ := json.Marshal(order)
	msg := &primitive.Message{
		Topic: "inventory_reback",
		Body:  jsonStr,
	}
	res, err := p.SendMessageInTransaction(context.Background(), msg)
	if err != nil {
		return nil, status.Error(codes.Internal, "sending msg fails")
	}
	if res.State == primitive.CommitMessageState {
		return nil, status.Error(ol.Code, ol.Detail)
	}
	return &proto.OrderInfoResponse{Id: ol.Id, OrderSn: order.OrderSn, Total: ol.OrderMount}, nil
}
func (*OrderServer) UpdateOrderStatus(c context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	if result := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "order not exists")
	}
	return &emptypb.Empty{}, nil
}
