package global

import (
	"gorm.io/gorm"
	"mxshop_srvs/stocks_srvs/config"
)

var (
	DB            *gorm.DB
	ServiceConfig = &config.ServiceConfig{}
	NacosConfig   = &config.NacosConfig{}
)
