package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"log"
	"time"

	"gorm.io/gorm"

	"xiaoyuzhou/pkg/setting"
)

var db *gorm.DB

//Model ...
type Model struct {
	ID        uint `gorm:"primarykey" json:"id"`
	CreatedAt int  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt int  `gorm:"column:updated_at" json:"updated_at"`
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

	//db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	//db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)

	// 用新函数替换GORM Create、Update流程自带的回调函数
	db.Callback().Create().Replace("gorm:before_create", updateTimeStampForBeforeCreateCallback)
	db.Callback().Update().Replace("gorm:before_update", updateTimeStampForBeforeUpdateCallback)

	db.AutoMigrate(&Author{}, &Article{}, &Category{}, &User{})
}

func updateTimeStampForBeforeCreateCallback(db *gorm.DB) {
	db.Statement.SetColumn("CreatedAt", time.Now().Unix())
}

func updateTimeStampForBeforeUpdateCallback(db *gorm.DB) {
	db.Statement.SetColumn("UpdatedAt", time.Now().Unix())
}
