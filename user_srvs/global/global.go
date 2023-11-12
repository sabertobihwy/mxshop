package global

import (
	"gorm.io/gorm"
	"mxshop_srvs/user_srvs/config"
)

var (
	DB            *gorm.DB
	ServiceConfig = &config.ServiceConfig{}
	NacosConfig   = &config.NacosConfig{}
)
