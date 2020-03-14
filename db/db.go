package db

import (
	"com.fs/event-service/config"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var gormDb *gorm.DB

// Open 打开数据库
func Open() {
	mysqlConfig := config.GetEventServiceConfig().DatabaseConfig
	template := "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	connStr := fmt.Sprintf(template, mysqlConfig.Username, mysqlConfig.Password, mysqlConfig.Address, mysqlConfig.Schema)
	openDb, err := gorm.Open("mysql", connStr)
	if err != nil {
		log.Println(err.Error())
		panic("数据库连接异常")
	}

	openDb.LogMode(true)

	gormDb = openDb

	autoMigrate(&Consumer{}, &ProducerPlugin{}, &Producer{})
}

// Close 关闭数据库
func Close() {
	_ = gormDb.Close()
	gormDb = nil
}

// GetInstance 获取数据库实例
func getInstance() *gorm.DB {
	return gormDb
}

func autoMigrate(values ...interface{}) {
	gormDb.AutoMigrate(values...)
}
