package order_service

import "xiaoyuzhou/models"

func CheckOrderIfPayed(order_id string) (bool, error) {
	return true, nil
}

// CreateOrderRecord: 创建订单记录
func CreateOrderRecord(OrderId string, PayMethod string, uid string, amount float32, tarotIdlist []int, ques string) error {
	return models.CreateOrder(OrderId, PayMethod, uid, amount, tarotIdlist, ques)
}
