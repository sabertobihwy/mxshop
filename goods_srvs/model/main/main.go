package main

import (
	"fmt"
	"log"
	"mxshop_srvs/goods_srvs/model"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := "root:root@tcp(192.168.2.112:3306)/mxshop_goods_srvs?charset=utf8mb4&parseTime=True&loc=Local"
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
		//	NamingStrategy: schema.NamingStrategy{TablePrefix: "mxshop_"},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	//DB.AutoMigrate(&model.Banner{}, &model.Brands{},
	//	&model.GoodsCategoryBrand{}, &model.Goods{}, &model.Category{})

	var brands = []model.Brands{}
	result := DB.Find(&brands)
	fmt.Println(int32(result.RowsAffected))

}
