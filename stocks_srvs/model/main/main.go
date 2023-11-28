package main

import (
	"fmt"
	"gorm.io/gorm/schema"
	"log"
	"mxshop_srvs/stocks_srvs/model"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := "root:root@tcp(192.168.2.112:3306)/mxshop_stocks?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			//IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			//ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful: true, // Disable color
		},
	)
	// NamingStrategy & TableName cannot config concurrently
	var err error
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix: "mxshop_",
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	DB.AutoMigrate(&model.InventoryDetail{})
	DB.Create(&model.InventoryDetail{
		OrderSn: "mimimi",
		Status:  1,
		Details: []model.GoodsInfo{
			{GoodsId: 1, Num: 2}, {GoodsId: 2, Num: 3}, {GoodsId: 3, Num: 4},
		},
	})
	var inv model.InventoryDetail
	DB.Where(&model.InventoryDetail{OrderSn: "mimimi"}).Find(&inv)
	for _, goodsInfo := range inv.Details {
		fmt.Println(goodsInfo.GoodsId, goodsInfo.Num)
	}

}
