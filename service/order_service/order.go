package order_service

import "xiaoyuzhou/models"

func CheckOrderIfPayed(OriOrderID, payMethod, uid string) (bool, error) {
	return models.CheckOrderIfPayed(OriOrderID, payMethod, uid)
}

// CreateOrderRecord: 创建订单记录
func CreateOrderRecord(OrderId string, OriOrderID string, PayMethod string, uid string, amount float32, tarotIdlist []int, ques string) error {
	return models.CreateOrder(OrderId, OriOrderID, PayMethod, uid, amount, tarotIdlist, ques)
}

func UpdateOrderStatus(OriOrderID string, status string, tansactionId string) error {
	return models.UpdateOrderStatus(OriOrderID, status, tansactionId)
}
