package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"xiaoyuzhou/pkg/logging"
)

type TrueWord struct {
	Model
	Word      string `gorm:"column:word;type:varchar(191);not null" json:"word"`                                         // 真言
	Language  string `gorm:"column:language;type:varchar(10);not null" json:"language" default:"jp" enums:"jp,zh,en,ts"` // 语言
	UpdatedBy string `gorm:"column:updated_by;type:varchar(10);not null" json:"updated_by"`
	CreatedBy string `gorm:"column:created_by;type:varchar(10);not null" json:"created_by"`
}

func AddTrueWord(lang string, wordL []string) error {
	toAdd := make([]TrueWord, 0)
	for _, value := range wordL {
		toAdd = append(toAdd, TrueWord{Language: lang, Word: value})
	}
	return Db.Create(&toAdd).Error
}

func TrueWordExistById(id int) bool {
	var tw TrueWord
	if err := Db.Model(&TrueWord{}).Where("id = ?", id).First(&tw).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			logging.Error(fmt.Sprintf("error: %s", err.Error()))
		}
		return false
	}
	return true
}

func EditTrueWord(id int, _data interface{}) error {
	if err := Db.Model(&TrueWord{}).Where("id = ?", id).Updates(_data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTrueWord(id int) error {
	// 先判断是否大于1
	var count int64
	var tw TrueWord
	if err := Db.Where("id = ?", id).First(&tw).Error; err != nil {
		return err
	}
	Db.Model(&TrueWord{}).Where("language = ?", tw.Language).Count(&count)
	if count > 1 {
		if err := Db.Where("id = ?", id).Delete(&TrueWord{}).Error; err != nil {
			return err
		}
		return nil
	}
	return errors.New("此类型真言只有一个了！")
}

func GetTrueWord(pageNum int, pageSize int, cond string, vals []interface{}) ([]TrueWord, int64, error) {
	var tws []TrueWord
	var count int64
	Db.Model(&TrueWord{}).Where(cond, vals...).Count(&count)
	err := Db.Where(cond, vals...).Order("created_at desc").Offset(pageNum).Limit(pageSize).Find(&tws).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	return tws, count, nil
}
