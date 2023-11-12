package global

import (
	"gorm.io/gorm"
	"mxshop_srvs/goods_srvs/config"
)

var (
	DB            *gorm.DB
	ServiceConfig = &config.ServiceConfig{}
	NacosConfig   = &config.NacosConfig{}
)
