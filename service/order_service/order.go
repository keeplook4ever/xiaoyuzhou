package order_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/paypal"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/logging"
)

const (
	PayPalClientIDTest = "AYxax5vAtOVCwWE7k3-IkfeeW91WssDy0X87o3hYU6v64yJPWDkb_9mWbcKrlixMDssEXnaU75Qwd8A3"
	PayPalSecretTest   = "EIIDNX57BlG8pCeFiiY6WJipeMtSQsIiOOP3ojg0_gtNSd3ndB0asdBQXP9IIeGh6gmlRnJHjjeozKGP"
)

func CheckOrderIfPayed(OrderID, uid string) (bool, error) {
	return models.CheckOrderIfPayed(OrderID, uid)
}

func GetOrderStatus(OrderID string) (int, error) {
	od, err := models.GetOrderByOrderId(OrderID)
	if err != nil {
		logging.Error(fmt.Sprintf("获取订单(%s)失败：%s", OrderID, err.Error()))
		return 0, err
	}
	oriOrderId := od.OriOrderID
	payMethod := od.PayMethod

	err = CaptureOrder(oriOrderId, payMethod)
	if err != nil {
		return od.Status, nil
	}

	od, err = models.GetOrderByOrderId(OrderID)
	if err != nil {
		return od.Status, err
	}
	return od.Status, nil
}

func AddPaymentInfoToOrder(OrderId string, OriOrderID string, PayMethod string, amount float32) error {
	return models.AddPaymentInfoToOrder(OrderId, OriOrderID, PayMethod, amount)
}

// CreateRecordWithNoOrder 创建没有订单的抽牌记录
func CreateRecordWithNoOrder(uid string, Question string, ts int64, tarotIdList []uint) (error, string) {
	return models.CreateRecordWithNoOrder(uid, Question, ts, tarotIdList)
}

func UpdateOrderStatus(OriOrderID string, payMethod string, status int, tansactionId string) error {
	return models.UpdateOrderStatus(OriOrderID, payMethod, status, tansactionId)
}

func CaptureOrder(OriOrderId, payMethod string) (err error) {
	logging.Debugf("Start Capture Order: %s", OriOrderId)
	// 检查订单是否已经支付或支付失败
	od, err := models.GetOrderByOriOrder(OriOrderId, payMethod)
	if err != nil {
		logging.Error(fmt.Sprintf("通过支付渠道OriOrderId(%s)获取订单失败：%s", OriOrderId, err.Error()))
	} else if od.Status != 0 {
		return errors.New("ok")
	}

	client, err := paypal.NewClient(PayPalClientIDTest, PayPalSecretTest, false)
	if err != nil {
		logging.Error(fmt.Sprintf("Error %v", err))
		return
	}
	// 打开Debug开关，输出日志
	client.DebugSwitch = gopay.DebugOn
	ctx := context.Background()
	ppRspc, err := client.OrderCapture(ctx, OriOrderId, nil)
	if err != nil {
		logging.Error(fmt.Sprintf("Error %v", err))
		return
	}
	status := 0

	if ppRspc.Code != paypal.Success {
		status = 2
		err = UpdateOrderStatus(ppRspc.Response.Id, payMethod, status, "")
		if err != nil {
			return errors.New(fmt.Sprintf("交易单 %s 支付失败, 订单状态更新失败！"))
		}
		logging.Debugf("Return Code %v", ppRspc.Code)
		return errors.New(fmt.Sprintf("Retuen Code is %d", ppRspc.Code))
	}

	// 判断订单状态 COMPLETED 代表完成
	transaction := ppRspc.Response.PurchaseUnits[0].Payments.Captures[0]
	orderIdReturn := ppRspc.Response.Id
	orderStatus := ppRspc.Response.Status
	if !(transaction.Status == "COMPLETED" && orderStatus == "COMPLETED" && orderIdReturn == OriOrderId) {
		logging.Error(fmt.Sprintf("交易单号：%s 失败，Status: %s", transaction.Id, transaction.Status))
		status = 2
		err = UpdateOrderStatus(orderIdReturn, payMethod, status, transaction.Id)
		if err != nil {
			return errors.New(fmt.Sprintf("交易单 %s 支付失败, 订单状态更新失败！", transaction.Id))
		}
		return errors.New(fmt.Sprintf("交易单 %s 支付失败", transaction.Id))
	} else {
		status = 1
		logging.Info(fmt.Sprintf("订单号：%s, 交易单号：%s, Status: %s", orderIdReturn, transaction.Id, orderStatus))
	}

	// 记录订单交易单号
	// 新订单状态为已支付，并在订单表中记录交易单号
	err = UpdateOrderStatus(orderIdReturn, payMethod, status, transaction.Id)
	if err != nil {
		logging.Error(fmt.Sprintf("更新订单状态失败：%v", err))
		return
	}
	logging.Info(fmt.Sprintf("订单 %s 更新成功", OriOrderId))
	return nil
}
