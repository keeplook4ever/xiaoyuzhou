package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"log"
	"time"

	"gorm.io/gorm"

	"xiaoyuzhou/pkg/setting"
)

var Db *gorm.DB

//Model ...
type Model struct {
	ID        uint `gorm:"primaryKey;not null;autoIncrement;type:int" json:"id"`
	CreatedAt int  `gorm:"column:created_at;not null;type:int" json:"created_at"`
	UpdatedAt int  `gorm:"column:updated_at;not null;type:int" json:"updated_at"`
}

// Setup initializes the database instance
func Setup() {
	var err error
	Db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name)))

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	sqlDB, err := Db.DB()
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	sqlDB.SetMaxIdleConns(100)          // 设置MySQL的最大空闲连接数（推荐100）
	sqlDB.SetMaxOpenConns(100)          // 设置MySQL的最大连接数（推荐100）
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置MySQL的空闲连接最大存活时间（推荐10s）

	// 用新函数替换GORM Create、Update流程自带的回调函数
	err = Db.Callback().Create().Replace("gorm:before_create", updateTimeStampForBeforeCreateCallback)
	if err != nil {
		log.Fatalf("models.Replace err: %v", err)
	}
	err = Db.Callback().Update().Replace("gorm:before_update", updateTimeStampForBeforeUpdateCallback)
	if err != nil {
		log.Fatalf("models.Replace err: %v", err)
	}

	err = Db.AutoMigrate(&Author{}, &Article{}, &Category{}, &User{}, &Lottery{}, &LuckyToday{}, &LotteryContent{})
	if err != nil {
		log.Fatalf("models.AutoMigrate err: %v", err)
	}
}

func updateTimeStampForBeforeCreateCallback(db *gorm.DB) {
	db.Statement.SetColumn("CreatedAt", time.Now().Unix())
}

func updateTimeStampForBeforeUpdateCallback(db *gorm.DB) {
	db.Statement.SetColumn("UpdatedAt", time.Now().Unix())
}
