package handler

import (
	"gorm.io/gorm"
	"mxshop_srvs/goods_srvs/proto"
)

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

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}
