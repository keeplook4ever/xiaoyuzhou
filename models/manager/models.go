package manager

import (
	"fmt"
	"gorm.io/driver/mysql"
	"log"

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

	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Author{}, &Article{}, &Category{}, &User{})
}
