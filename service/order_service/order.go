package order_service

import "xiaoyuzhou/models"

func CheckOrderIfPayed(orderId string) (bool, error) {
	return models.CheckOrderIfPayed(orderId)
}

// CreateOrderRecord: 创建订单记录
func CreateOrderRecord(OrderId string, PayMethod string, uid string, amount float32, tarotIdlist []int, ques string) error {
	return models.CreateOrder(OrderId, PayMethod, uid, amount, tarotIdlist, ques)
}

func UpdateOrderStatus(OrderId string, status string, tansactionId string) error {
	return models.UpdateOrderStatus(OrderId, status, tansactionId)
}
