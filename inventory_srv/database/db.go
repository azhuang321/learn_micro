package database

import (
	"fmt"
	"inventory_srv/global"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func myLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)
}

func GetDB() (db *gorm.DB, err error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		global.ServerConfig.Mysql.User,
		global.ServerConfig.Mysql.Password,
		global.ServerConfig.Mysql.Host,
		global.ServerConfig.Mysql.Port,
		global.ServerConfig.Mysql.Db,
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: myLogger(),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
