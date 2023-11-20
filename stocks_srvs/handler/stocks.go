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
	//global.DB.First(&inv, g.GoodsId)
	inv.Goods = g.GoodsId
	inv.Stocks = g.Num
	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}
func (s *StocksServer) InvDetail(c context.Context, g *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.First(&inv, g.GoodsId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "inventory not found")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil

}
