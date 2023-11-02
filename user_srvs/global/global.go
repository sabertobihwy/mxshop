package global

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	DB *gorm.DB
)

func init() {
	dsn := "root:root@tcp(192.168.2.112:3306)/mxshop?charset=utf8mb4&parseTime=True&loc=Local"
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
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{TablePrefix: "mxshop_"},
		Logger:         newLogger,
	})
	if err != nil {
		panic(err)
	}

}
