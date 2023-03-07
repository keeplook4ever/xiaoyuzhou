package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"math/rand"
	"time"
	"xiaoyuzhou/pkg/util"
)

type Order struct {
	Model
	OrderId   string  `gorm:"column:order_id;not null;type:varchar(191)" json:"order_id"`                                         // 订单ID
	Uid       string  `gorm:"column:uid;not null;type:varchar(191)" json:"uid"`                                                   // 订单付款用户
	Status    int     `gorm:"column:status;not null;type:tinyint(4);default:0" json:"status" enums:"0,1"`                         // 订单状态 0:未支付，1：已支付
	Amount    float32 `gorm:"column:amount;not null;type:float" json:"amount"`                                                    // 订单实际支付金额
	TarotList string  `gorm:"column:tarot_list;not null;type:varchar(100)" json:"tarot_list"`                                     // 塔罗牌id列表：单张是1个，三张是3个
	PayedTime int     `gorm:"column:payed_time;not null;type:tinyint(20);default:0" json:"payed_time"`                            // 支付时间戳
	PayMethod string  `gorm:"column:pay_method;not null;type:varchar(100)" json:"pay_method" enums:"paypal,wechat,alipay,credit"` // 支付方式：PayPal,微信,支付宝,信用卡
	Ques      string  `gorm:"column:ques;not null;type:varchar(191)" json:"ques"`                                                 // 用户输入的问题
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
func GetOneTarotFromOrder(orderId string) (*TarotDto, string, error) {
	var tarot Tarot
	var od Order
	if err := Db.Model(&Order{}).Where("order_id = ?", orderId).Find(&od).Error; err != nil {
		return nil, "", err
	}
	tLSlice := util.StringToIntSlice(od.TarotList)

	if err := Db.Model(&Tarot{}).Where("id in ?", tLSlice).Find(&tarot).Error; err != nil {
		return nil, "", err
	}
	resp := tarot.ToTarotDto()

	// 将answerList的值变成其中之一

	rand.Seed(time.Now().Unix())
	answers := make([]string, 0)
	answers = append(answers, resp.AnswerList[rand.Intn(len(resp.AnswerList))])
	resp.AnswerList = answers
	return &resp, od.Ques, nil
}

func CreateOrder(OrderId string, PayMethod string, uid string, amount float32, tarotIdlist []int, question string) error {
	ids, err := json.Marshal(tarotIdlist)
	if err != nil {
		return err
	}
	if err := Db.Create(&Order{
		OrderId:   OrderId,
		Uid:       uid,
		Status:    0,
		Amount:    amount,
		TarotList: string(ids),
		PayMethod: PayMethod,
		Ques:      question,
	}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateOrderStatus(OrderId string, status string) error {
	data := map[string]interface{}{
		"status":     status,
		"payed_time": time.Now().UnixMilli(),
	}
	if err := Db.Model(&Order{}).Where("order_id = ?", OrderId).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
