package main

import (
	"crypto/sha512"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"mxshop_srvs/user_srvs/model"
)

func main() {
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
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{TablePrefix: "mxshop_"},
		Logger:         newLogger,
	})
	if err != nil {
		panic(err)
	}

	options := &password.Options{16, 100, 30, sha512.New} // use sha512 better than MD5
	salt, encodedPwd := password.Encode("admin123", options)
	newPwd := fmt.Sprintf("$sha512$%s$%s", salt, encodedPwd)

	for i := 0; i < 10; i++ {
		user := model.User{
			NickName: fmt.Sprintf("bobby%d", i),
			Mobile:   fmt.Sprintf("1996209957%d", i),
			Password: newPwd,
		}
		DB.Save(&user)
	}

}
