package order_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/paypal"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/pkg/setting"
)

const (
	PayPalClientIDTest = "Adume2QcPr8y_RdWg2aVXnukqXC1zQWmOJmgN_7F-fMNr-LOHPlmBypWSz2QVXqzVw9Pzx18mPnwqVId"
	PayPalSecretTest   = "EJPUmppFw_NP36NBxZUO3wLoWu3MkKceVd7X8rosU7LMCUDzNB-XTl5bVtaYNzteGq6DhC9dEFGEUfvb"
	PayPalClientIDPrd  = "AZTKeoDxoleGB9frtSlHugzdsM1kLcf-VsycOgq2lkGVquO4lOW5f43asbB32n1yBxFZ1xglY363D3hk"
	PayPalSecretPrd    = "EK_ulMtKhhSDRXdU2C7r6XrgEpjywwtvrLOGcDEdr8bwbejuyx3slE7QbGsDraoSLamAezBZ7fq32MwF"
)

func CheckOrderIfPayed(OrderID string) (bool, error) {
	return models.CheckOrderIfPayed(OrderID)
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
func CreateRecordWithNoOrder(uid string, Question string, isMobile bool, ts int64, tarotIdList []uint) (error, string) {
	return models.CreateRecordWithNoOrder(uid, Question, isMobile, ts, tarotIdList)
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

	client, err := paypal.NewClient(PayPalClientIDPrd, PayPalSecretPrd, true)
	if err != nil {
		logging.Error(fmt.Sprintf("创建clientPrd失败 %v", err))
		return
	}
	if setting.PaymentSetting.Mode == "debug" {
		client, err = paypal.NewClient(PayPalClientIDTest, PayPalSecretTest, false)
		if err != nil {
			logging.Error(fmt.Sprintf("创建client测试失败 %v", err))
			return
		}
	}

	// 打开Debug开关，输出日志
	client.DebugSwitch = gopay.DebugOn
	ctx := context.Background()
	ppRspc, err := client.OrderCapture(ctx, OriOrderId, nil)
	if err != nil {
		logging.Error(fmt.Sprintf("捕获订单状态失败： %v", err))
		return
	}
	status := 0

	if ppRspc.Code != paypal.Success {
		status = 2
		err = UpdateOrderStatus(ppRspc.Response.Id, payMethod, status, "")
		if err != nil {
			logging.Error(fmt.Sprintf("更新订单状态失败: %s", err.Error()))
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
			logging.Error(fmt.Sprintf("更新订单状态失败: %s", err.Error()))
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
