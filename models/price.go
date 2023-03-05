package models

type Price struct {
	Model
	SingleOrig       float32 `gorm:"column:single_orig;not null;type:float" json:"single_orig"`               // 单个原价
	SingleSellHigher float32 `gorm:"column:single_sell_higher;not null;type:float" json:"single_sell_higher"` // 单个较高售价
	SingleSellLower  float32 `gorm:"column:single_sell_lower;not null;type:float" json:"single_sell_lower"`   // 单个较低售价
	ThreeOrig        float32 `gorm:"column:three_orig;not null;type:float" json:"three_orig"`                 // 三个原价
	ThreeSellHigher  float32 `gorm:"column:three_sell_higher;not null;type:float" json:"three_sell_higher"`   // 三个较高售价
	ThreeSellLower   float32 `gorm:"column:three_sell_lower;not null;type:float" json:"three_sell_lower"`     // 三个较低售价
	CreatedBy        string  `gorm:"column:created_by;not null;type:varchar(50)" json:"created_by"`           // 创建者
	UpdatedBy        string  `gorm:"column:updated_by;not null;type:varchar(50)" json:"updated_by"`           // 更新者
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
	}
	var has []Price
	// 有的话先删除
	if err := Db.Model(&Price{}).Find(&has).Error; err != nil {
		return err
	} else {
		Db.Where("1 = 1").Delete(&Price{})
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
	if err := Db.Model(&Price{}).Where("id > ?", 0).Updates(&setD).Error; err != nil {
		return err
	}
	return nil
}

func GetPrice() (*Price, error) {
	var res Price
	if err := Db.Model(&Price{}).First(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}
