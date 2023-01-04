package manager

import (
	"fmt"
	"gorm.io/driver/mysql"
	"log"
	"time"

	"gorm.io/gorm"

	"xiaoyuzhou/pkg/setting"
)

var db *gorm.DB

type Model struct {
}

// Setup initializes the database instance
func Setup() {
	var err error
	db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name)))

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(100)          // 设置MySQL的最大空闲连接数（推荐100）
	sqlDB.SetMaxOpenConns(100)          // 设置MySQL的最大连接数（推荐100）
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置MySQL的空闲连接最大存活时间（推荐10s）
	db.AutoMigrate(&Author{}, &Article{}, &Category{}, &User{})
}
