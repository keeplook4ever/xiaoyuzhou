package models

import "xiaoyuzhou/pkg/logging"

type Price struct {
	Model
	SingleOrig       float32 `gorm:"column:single_orig;not null;type:float" json:"single_orig"`                     // 单个原价
	SingleSellHigher float32 `gorm:"column:single_sell_higher;not null;type:float" json:"single_sell_higher"`       // 单个较高售价
	SingleSellLower  float32 `gorm:"column:single_sell_lower;not null;type:float" json:"single_sell_lower"`         // 单个较低售价
	ThreeOrig        float32 `gorm:"column:three_orig;not null;type:float" json:"three_orig"`                       // 三个原价
	ThreeSellHigher  float32 `gorm:"column:three_sell_higher;not null;type:float" json:"three_sell_higher"`         // 三个较高售价
	ThreeSellLower   float32 `gorm:"column:three_sell_lower;not null;type:float" json:"three_sell_lower"`           // 三个较低售价
	Language         string  `gorm:"column:language;not null;type:varchar(10)" json:"language" enums:"jp,zh,en,tc"` // 语言
	CreatedBy        string  `gorm:"column:created_by;not null;type:varchar(50)" json:"created_by"`                 // 创建者
	UpdatedBy        string  `gorm:"column:updated_by;not null;type:varchar(50)" json:"updated_by"`                 // 更新者
}

func SetPrice(data map[string]interface{}) error {
	setD := Price{
		SingleOrig:       data["single_orig"].(float32),
		SingleSellHigher: data["single_sell_higher"].(float32),
		SingleSellLower:  data["single_sell_lower"].(float32),
		ThreeOrig:        data["three_orig"].(float32),
		ThreeSellHigher:  data["three_sell_higher"].(float32),
		ThreeSellLower:   data["three_sell_lower"].(float32),
		CreatedBy:        data["created_by"].(string), // 创建者
		UpdatedBy:        data["updated_by"].(string), // 更新者
		Language:         data["language"].(string),   // 地区
	}
	var has []Price
	// 有的话先删除
	if err := Db.Model(&Price{}).Where("language = ?", data["language"]).Find(&has).Error; err != nil {
		return err
	} else {
		Db.Where("language = ?", data["language"]).Delete(&Price{})
	}

	if err := Db.Model(&Price{}).Create(&setD).Error; err != nil {
		return err
	}
	return nil
}

func UpdatePrice(data map[string]interface{}) error {
	setD := Price{
		SingleOrig:       data["single_orig"].(float32),
		SingleSellHigher: data["single_sell_higher"].(float32),
		SingleSellLower:  data["single_sell_lower"].(float32),
		ThreeOrig:        data["three_orig"].(float32),
		ThreeSellHigher:  data["three_sell_higher"].(float32),
		ThreeSellLower:   data["three_sell_lower"].(float32),
	}
	// 如果是0则updates会自动忽略更新
	if err := Db.Model(&Price{}).Where("language = ?", data["language"]).Updates(&setD).Error; err != nil {
		return err
	}
	return nil
}

func GetPrice(language string) (*Price, error) {
	var res Price
	if err := Db.Model(&Price{}).Where("language = ?", language).First(&res).Error; err != nil {
		logging.Debugf("Error %s", err.Error())
		return nil, err
	}
	return &res, nil
}

func GetPriceTotal() ([]Price, error) {
	var res []Price
	if err := Db.Model(&Price{}).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func GetPaymentPrice(scene, language string) float32 {
	// enums:"ta_one_high,ta_one_low,ta_three_high,ta_three_low"
	// language: jp,zh,en,tc
	priceTotal, err := GetPrice(language)
	if err != nil {
		if language == "jp" {
			return 2480
		} else if language == "zh" {
			return 98
		} else if language == "tc" {
			return 598
		} else if language == "en" {
			return 24.99
		}
	}

	switch scene {
	case "ta_one_high":
		return priceTotal.SingleSellHigher
	case "ta_one_low":
		return priceTotal.SingleSellLower
	case "ta_three_high":
		return priceTotal.ThreeSellHigher
	case "ta_three_low":
		return priceTotal.ThreeSellLower
	}
	return priceTotal.SingleSellHigher // 默认价格
}
