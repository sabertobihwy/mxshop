package global

import (
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	"mxshop_srvs/order_srvs/config"
	"mxshop_srvs/order_srvs/proto"
)

var (
	DB            *gorm.DB
	ServiceConfig = &config.ServiceConfig{}
	NacosConfig   = &config.NacosConfig{}
	Redsync       *redsync.Redsync
	GoodsClient   proto.GoodsClient
	StocksClient  proto.StocksClient
)
