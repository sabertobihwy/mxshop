package global

import (
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	"mxshop_srvs/stocks_srvs/config"
)

var (
	DB            *gorm.DB
	ServiceConfig = &config.ServiceConfig{}
	NacosConfig   = &config.NacosConfig{}
	Redsync       *redsync.Redsync
)
