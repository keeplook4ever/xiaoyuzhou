package models

import (
	"gorm.io/gorm"
	"xiaoyuzhou/pkg/util"
)

type Order struct {
	Model
	OrderId   string  `gorm:"column:order_id;not null;type:varchar(191)" json:"order_id"`                   // 订单ID
	Uid       string  `gorm:"column:uid;not null;type:varchar(191)" json:"uid"`                             // 订单付款用户
	Status    int     `gorm:"column:status;not null;type:tinyint(4);default:0" json:"status" enums:"0,1,2"` // 订单状态 0:刚创建，未支付，1：已支付成功，2：支付
	Amount    float32 `gorm:"column:amount;not null;type:float" json:"amount"`                              // 订单实际支付金额
	TarotList string  `gorm:"column:tarot_list;not null;type:varchar(100)" json:"tarot_list"`               // 塔罗牌id列表
	PayedTime int     `gorm:"column:payed_time;not null;type:int;default:0" json:"payed_time"`              // 支付时间戳

}

func CheckOrderIfValid(orderId, uid string) (bool, error) {
	var od Order
	if err := Db.Model(&Order{}).Where("uid = ? and order_id = ?", uid, orderId).First(&od).Error; err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	} else if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if od.ID > 0 {
		return true, nil
	}
	return false, nil
}

// GetOneTarotFromOrder 根据订单号获取对应的塔罗牌
func GetOneTarotFromOrder(orderId string) (*TarotDto, error) {
	var tarot Tarot
	var tIdList string
	if err := Db.Model(&Order{}).Where("order_id = ?", orderId).Pluck("tarot_list", tIdList).Error; err != nil {
		return nil, err
	}
	tLSlice := util.StringToIntSlice(tIdList)

	if err := Db.Model(&Tarot{}).Where("id in ?", tLSlice).Find(&tarot).Error; err != nil {
		return nil, err
	}
	resp := tarot.ToTarotDto()
	return &resp, nil
}
