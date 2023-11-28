package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"mxshop_srvs/stocks_srvs/global"
	"mxshop_srvs/stocks_srvs/model"
	"mxshop_srvs/stocks_srvs/proto"
	"sync"
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

var l sync.Mutex

// reduce inventory
// 本地事务保证 input[1:5, 2:10, 3:20]全部执行，不能保证多线程下数据一致
func (s *StocksServer) Sell(c context.Context, sell *proto.SellInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin() // start local transaction
	var goodsList model.GooodsInfoList
	for _, g := range sell.GoodsInvInfo {
		var inv model.Inventory
		mutex := global.Redsync.NewMutex(fmt.Sprintf("goods_%d", g.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "redsync lock error")
		}
		if result := global.DB.
			Where(&model.Inventory{Goods: g.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // rollback
			return nil, status.Errorf(codes.InvalidArgument, "inventory not found")
		}
		if inv.Stocks < g.Num {
			tx.Rollback() // rollback
			return nil, status.Errorf(codes.ResourceExhausted, "inventory not enough")
		}
		inv.Stocks -= g.Num // 数据一致要靠分布式锁
		tx.Save(&inv)
		goodsList = append(goodsList, model.GoodsInfo{
			GoodsId: g.GoodsId,
			Num:     g.Num,
		})
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "redsync unlock error")
		}
		// upate inventory where goods = inv.goodid, version = inv.version set stocks = inv.Stocks- g.num and version = inv.version+1
		//if result := global.DB.Model(&model.Inventory{}).Select("stocks", "verson").
		//	Where("goods = ? and verson = ?", inv.Goods, inv.Verson).
		//	Updates(model.Inventory{Stocks: inv.Stocks - g.Num, Verson: inv.Verson + 1}); result.RowsAffected == 0 {
		//	zap.S().Info("reducing stocks failed")
		//} else {
		//	break
		//}

	}
	if result := tx.Create(&model.InventoryHistory{
		OrderSn: sell.OrderSn,
		Status:  1,
		Details: goodsList,
	}); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "saving stocks history fails")
	}
	tx.Commit() // commit
	return &emptypb.Empty{}, nil
}

// stock-return situation:
// 1. timeout 2. refund 3. order failed
func (s *StocksServer) Reback(c context.Context, sell *proto.SellInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin() // start local transaction
	l.Lock()
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
	l.Unlock()
	return &emptypb.Empty{}, nil
}

func AutoReback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderSn string `gorm:"type:varchar(30);index"`
	}
	var oi OrderInfo
	for _, msg := range msgs {
		if err := json.Unmarshal(msg.Body, &oi); err != nil {
			zap.S().Debugf("JSON analysing error")
			return consumer.ConsumeSuccess, nil
		}
		var inv model.InventoryHistory
		if result := global.DB.Where(&model.InventoryHistory{OrderSn: oi.OrderSn, Status: 1}).First(&inv); result.RowsAffected == 0 {
			// already reback
			return consumer.ConsumeSuccess, nil
		}
		tx := global.DB.Begin()
		for _, goodsInfo := range inv.Details {
			if result := tx.Where(&model.Inventory{Goods: goodsInfo.GoodsId}).Update("stocks", gorm.Expr("stocks+?", goodsInfo.Num)); result.RowsAffected == 0 {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}
		if result := tx.Where(&model.InventoryHistory{OrderSn: oi.OrderSn}).Update("status", 2); result.RowsAffected == 0 {
			tx.Rollback()
			return consumer.ConsumeRetryLater, nil
		}
		tx.Commit()
	}
	return consumer.ConsumeSuccess, nil
}
