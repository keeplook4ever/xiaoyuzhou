package models

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/pkg/util"
)

type Order struct {
	Model
	OrderId       string  `gorm:"column:order_id;not null;type:varchar(191)" json:"order_id"`                                         // 订单ID: 抽塔罗牌记录唯一id: uid+timestamp
	OriOrderID    string  `gorm:"column:ori_order_id;not null;type:varchar(191)" json:"ori_order_id"`                                 // 原始平台订单id, 比如paypal返回的id
	Uid           string  `gorm:"column:uid;not null;type:varchar(191)" json:"uid"`                                                   // 订单付款用户
	Amount        float32 `gorm:"column:amount;not null;type:float" json:"amount"`                                                    // 订单实际支付金额
	Status        int     `gorm:"column:status;not null;type:tinyint(3);default:0" json:"status"`                                     // 订单状态：0未支付，1已支付，2支付失败
	TarotList     string  `gorm:"column:tarot_list;not null;type:varchar(100)" json:"tarot_list"`                                     // 塔罗牌id列表：单张是1个，三张是3个
	PickTime      int64   `gorm:"column:pick_time;not null;type:int;default:0" json:"pick_time"`                                      // 塔罗抽取时间
	PayedTime     int     `gorm:"column:payed_time;not null;type:int;default:0" json:"payed_time"`                                    // 支付时间戳
	PayMethod     string  `gorm:"column:pay_method;not null;type:varchar(100)" json:"pay_method" enums:"paypal,wechat,alipay,credit"` // 支付方式：PayPal,微信,支付宝,信用卡
	Ques          string  `gorm:"column:ques;not null;type:varchar(191)" json:"ques"`                                                 // 用户输入的问题
	TransactionId string  `gorm:"column:transaction_id;not null;type:varchar(190)" json:"transaction_id"`                             // 交易付款流水号
}

// GetOneTarotFromOrder 根据订单号获取对应的塔罗牌
func GetOneTarotFromOrder(OrderId, uid string) (*TarotDto, string, int64, error) {
	var tarot Tarot
	var od Order
	if err := Db.Model(&Order{}).Where("order_id = ? and uid= ?", OrderId, uid).Find(&od).Error; err != nil {
		return nil, "", 0, err
	}
	tLSlice := util.StringToIntSlice(od.TarotList)

	if err := Db.Model(&Tarot{}).Where("id in ?", tLSlice).Find(&tarot).Error; err != nil {
		return nil, "", 0, err
	}
	resp := tarot.ToTarotDto()

	// 将answerList的值变成其中之一
	rand.Seed(time.Now().Unix())
	answers := make([]string, 0)
	answers = append(answers, resp.AnswerList[rand.Intn(len(resp.AnswerList))])
	resp.AnswerList = answers
	return &resp, od.Ques, od.PickTime, nil
}

func AddPaymentInfoToOrder(OrderId string, OriOrderID string, PayMethod string, amount float32) error {
	var od Order
	if err := Db.Model(&Order{}).Where("order_id = ?", OrderId).First(&od).Error; err != nil {
		return err
	}
	data := map[string]interface{}{
		"ori_order_id": OriOrderID,
		"pay_method":   PayMethod,
		"amount":       amount,
		"status":       0,
	}
	if err := Db.Model(&Order{}).Where("order_id = ?", OrderId).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// CreateRecordWithNoOrder 用户抽取塔罗牌输入问题，记录抽取时间
func CreateRecordWithNoOrder(uid, question string, ts int64, tarotIdlist []uint) (error, string) {
	ids, err := json.Marshal(tarotIdlist)
	if err != nil {
		return err, ""
	}
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	day := time.Now().Format("02")
	ls := strconv.FormatInt(ts, 10)[5:]

	// 订单号生成："TA" + 年后两位+月日 +秒的后5位 + uid后3位
	orderId := "TA" + year[2:] + month + day + ls + uid[33:]

	if err = Db.Create(&Order{
		OrderId:   orderId,
		Uid:       uid,
		Ques:      question,
		PickTime:  ts,
		TarotList: string(ids),
	}).Error; err != nil {
		return err, ""
	}
	return nil, orderId
}

func UpdateOrderStatus(OriOrderId string, payMethod string, status int, tansActionId string) error {
	var od Order
	data := map[string]interface{}{
		"status":         status,
		"payed_time":     time.Now().UnixMilli(),
		"transaction_id": tansActionId,
	}
	if err := Db.Model(&Order{}).Where("ori_order_id = ? and pay_method = ?", OriOrderId, payMethod).First(&od).Error; err == gorm.ErrRecordNotFound {
		logging.Debugf("Error: %s", err.Error())
		return err
	}

	if err := Db.Model(&Order{}).Where("ori_order_id = ? and pay_method = ?", OriOrderId, payMethod).Updates(data).Error; err != nil {
		logging.Debugf("Error: %s", err.Error())
		return err
	}
	return nil
}

func CheckOrderIfPayed(orderId, uid string) (bool, error) {
	var od Order
	if err := Db.Model(&Order{}).Where("order_id = ? and uid = ?", orderId, uid).First(&od).Error; err != nil {
		return false, err
	}

	// Status 改成0，1，2 : 0未支付，1已支付，2支付失败
	if od.Status == 1 && od.TransactionId != "" && od.PayedTime != 0 {
		return true, nil
	} else if od.Status == 2 {
		return false, errors.New("支付失败")
	}
	return false, nil

}

func GetOrderByOrderId(orderId string) (*Order, error) {
	var od Order
	if err := Db.Model(&Order{}).Where("order_id = ?", orderId).First(&od).Error; err != nil {
		return nil, err
	}
	return &od, nil
}

func GetOrderByOriOrder(OriOrderID, payMethod string) (*Order, error) {
	var od Order
	if err := Db.Model(&Order{}).Where("ori_order_id = ? and pay_method = ?", OriOrderID, payMethod).First(&od).Error; err != nil {
		return nil, nil
	}
	return &od, nil
}
