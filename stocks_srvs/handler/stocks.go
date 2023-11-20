package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/stocks_srvs/global"
	"mxshop_srvs/stocks_srvs/model"
	"mxshop_srvs/stocks_srvs/proto"
)

type StocksServer struct {
	proto.UnimplementedStocksServer
}

func (s *StocksServer) SetInv(c context.Context, g *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: g.GoodsId}).First(&inv) // Get the ID in &inv to update
	inv.Goods = g.GoodsId
	inv.Stocks = g.Num
	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}
func (s *StocksServer) InvDetail(c context.Context, g *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: g.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "inventory not found")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil

}

// reduce inventory
// 本地事务保证 input[1:5, 2:10, 3:20]全部执行，不能保证多线程下数据一致
func (s *StocksServer) Sell(c context.Context, sell *proto.SellInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin() // start local transaction
	for _, g := range sell.GoodsInvInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: g.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // rollback
			return nil, status.Errorf(codes.InvalidArgument, "inventory not found")
		}
		if inv.Stocks < g.Num {
			tx.Rollback() // rollback
			return nil, status.Errorf(codes.ResourceExhausted, "inventory not enough")
		}
		inv.Stocks -= g.Num // 数据一致要靠分布式锁
		tx.Save(&inv)
	}
	tx.Commit() // commit
	return &emptypb.Empty{}, nil
}

// stock-return situation:
// 1. timeout 2. refund 3. order failed
func (s *StocksServer) Reback(c context.Context, sell *proto.SellInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin() // start local transaction
	for _, g := range sell.GoodsInvInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: g.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // rollback
			return nil, status.Errorf(codes.InvalidArgument, "inventory not found")
		}

		inv.Stocks += g.Num // 数据一致要靠分布式锁
		tx.Save(&inv)
	}
	tx.Commit() // commit
	return &emptypb.Empty{}, nil
}
