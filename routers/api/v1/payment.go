package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/paypal"
	"github.com/go-pay/gopay/pkg/xlog"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/service/order_service"
	"xiaoyuzhou/service/tarot_service"
)

const (
	PayPalClientID = "AYxax5vAtOVCwWE7k3-IkfeeW91WssDy0X87o3hYU6v64yJPWDkb_9mWbcKrlixMDssEXnaU75Qwd8A3"
	PayPalSecret   = "EIIDNX57BlG8pCeFiiY6WJipeMtSQsIiOOP3ojg0_gtNSd3ndB0asdBQXP9IIeGh6gmlRnJHjjeozKGP"
)

type PayPalOrderResp struct {
	OrderID string `json:"order_id"`
}

type CreatePayPalOrderForm struct {
	CardType      string `json:"card_type" binding:"required" enums:"one,three"`      // 卡牌类型 one:单张，three: 多张
	Uid           string `json:"uid" binding:"required"`                              // 用户ID
	ReturnURL     string `json:"return_url" binding:"required"`                       // 支付成功返回URL
	CancelURL     string `json:"cancel_url"  binding:"required"`                      // 取消支付URL
	TarotIdList   []int  `json:"tarot_id_list" binding:"required"`                    // 塔罗牌id列表
	HigherOrLower string `json:"higher_or_lower" binding:"required" enums:"high,low"` // 高价格还是低价格
	Question      string `json:"question" binding:"required"`                         // 用户问题
}

// CreatePayPalOrder
// @Summary 创建PapPal支付订单
// @Accept json
// @Produce  json
// @Param _ body CreatePayPalOrderForm true "创建订单请求参数"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /player/paypal/checkout/orders [post]
// @Tags Player
func CreatePayPalOrder(c *gin.Context) {
	appG := app.Gin{C: c}
	var formD CreatePayPalOrderForm
	if err := c.ShouldBindJSON(&formD); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	client, err := paypal.NewClient(PayPalClientID, PayPalSecret, false)
	if err != nil {
		xlog.Error(err)
		appG.Response(http.StatusOK, "初始化client失败", nil)
		return
	}
	// 打开Debug开关，输出日志
	client.DebugSwitch = gopay.DebugOn

	xlog.Debugf("Appid: %s", client.Appid)
	xlog.Debugf("AccessToken: %s", client.AccessToken)
	xlog.Debugf("ExpiresIn: %d", client.ExpiresIn)
	xlog.Debugf("Uid: %s", formD.Uid)
	xlog.Debugf("ReturnURL: %s", formD.ReturnURL)
	xlog.Debugf("CacelURL: %s", formD.CancelURL)
	// Create Orders example
	var pus []*paypal.PurchaseUnit

	// 获取支付价格
	amount := tarot_service.GetPaymentPrice(formD.CardType, formD.HigherOrLower)
	var item = &paypal.PurchaseUnit{
		//ReferenceId: "TX12333331231232",
		//TODO:
		Amount: &paypal.Amount{
			CurrencyCode: "USD",
			// 金额从数据库查, 不支持小数
			Value: fmt.Sprintf("%.2f", amount),
		},
	}
	pus = append(pus, item)

	bm := make(gopay.BodyMap)
	bm.Set("intent", "CAPTURE").
		Set("purchase_units", pus).
		SetBodyMap("application_context", func(b gopay.BodyMap) {
			b.Set("brand_name", "小小の宇宙").
				//Set("locale", "en-PT").
				Set("return_url", formD.ReturnURL).
				Set("cancel_url", formD.CancelURL)
		})
	ctx := context.Background()
	ppRsp, err := client.CreateOrder(ctx, bm)
	if err != nil {
		xlog.Error(err)
		appG.Response(http.StatusOK, "创建订单失败", nil)
		return
	}
	if ppRsp.Code != paypal.Success {
		// do something
		xlog.Error("!!!")
		appG.Response(http.StatusOK, "创建订单失败", nil)
		return
	}

	xlog.Debugf("Response: %v", ppRsp.Response)

	// 13位时间戳
	ts13 := int(time.Now().UnixMilli())
	ts := strconv.Itoa(ts13)
	PayPalOrderId := "TA" + "-" + ppRsp.Response.Id + "-" + ts

	//将订单落库
	err = order_service.CreateOrderRecord(PayPalOrderId, "paypal", formD.Uid, amount, formD.TarotIdList, formD.Question)
	if err != nil {
		appG.Response(http.StatusOK, "订单落库失败", nil)
		return
	}

	xlog.Debugf("OrderID: %s", PayPalOrderId)
	appG.Response(http.StatusOK, e.SUCCESS, ppRsp.Response)
}

// ConfirmPayment
// @Summary 确认支付
// @Accept json
// @Produce  json
// @Param order_id path string true "订单号"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /player/paypal/confirm/orders/{order_id} [post]
// @Tags Player
func ConfirmPayment(c *gin.Context) {
	appG := app.Gin{
		C: c,
	}
	orderId := c.Param("order_id")
	client, err := paypal.NewClient(PayPalClientID, PayPalSecret, false)
	if err != nil {
		xlog.Error(err)
		appG.Response(http.StatusOK, "初始化client失败", nil)
		return
	}
	// 打开Debug开关，输出日志
	client.DebugSwitch = gopay.DebugOn
	ctx := context.Background()
	ppRspc, err := client.OrderConfirm(ctx, orderId, nil)
	if err != nil || ppRspc.Code != paypal.Success {
		xlog.Error(err)
		appG.Response(http.StatusOK, "支付订单失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, ppRspc.Response)
}

// CapturePayPalOrder
// @Summary 捕获PapPal支付订单
// @Accept json
// @Produce  json
// @Param order_id path string true "订单号"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /player/paypal/capture/orders/{order_id} [post]
// @Tags Player
//func CapturePayPalOrder(c *gin.Context) {
//	appG := app.Gin{C: c}
//	orderId := c.Param("order_id")
//	client, err := paypal.NewClient(PayPalClientID, PayPalSecret, false)
//	if err != nil {
//		xlog.Error(err)
//		appG.Response(http.StatusOK, "初始化client失败", nil)
//		return
//	}
//	// 打开Debug开关，输出日志
//	client.DebugSwitch = gopay.DebugOn
//	ctx := context.Background()
//	ppRspc, err := client.OrderCapture(ctx, orderId, nil)
//	if err != nil {
//		xlog.Error(err)
//		appG.Response(http.StatusOK, "捕获订单失败", nil)
//		return
//	}
//	if ppRspc.Code != paypal.Success {
//		// TODO ？？
//
//		appG.Response(http.StatusOK, "捕获订单失败", nil)
//		return
//	}
//
//	// TODO：判断订单状态 COMPLETED 代表完成
//
//	transaction := ppRspc.Response.PurchaseUnits[0].Payments.Captures[0]
//	if transaction.Status == "COMPLETED" {
//		xlog.Debugf("交易单号：%s", transaction.Id)
//	}
//	if transaction.Id != orderId {
//		xlog.Debugf("------")
//		return
//	}
//	xlog.Debug("Bingoo!")
//	appG.Response(http.StatusOK, e.SUCCESS, ppRspc.Response)
//}

// GetPayPalOrderDetail
// @Summary 获取PayPal订单详情
// @Accept json
// @Produce  json
// @Param order_id path string true "订单号"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /player/paypal/checkout/orders/{order_id} [get]
// @Tags Player
func GetPayPalOrderDetail(c *gin.Context) {
	appG := app.Gin{C: c}
	orderId := c.Param("order_id")
	client, err := paypal.NewClient(PayPalClientID, PayPalSecret, false)
	if err != nil {
		xlog.Error(err)
		appG.Response(http.StatusOK, "初始化client失败", nil)
		return
	}
	// 打开Debug开关，输出日志
	client.DebugSwitch = gopay.DebugOn
	ctx := context.Background()
	ppRspc, err := client.OrderDetail(ctx, orderId, nil)
	if err != nil {
		xlog.Error(err)
		appG.Response(http.StatusOK, "获取订单详情失败", nil)
		return
	}

	if ppRspc.Code != paypal.Success {
		// TODO ？？

		appG.Response(http.StatusOK, "获取订单详情失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, ppRspc.Response)

}

// ReceiveOrderEventsFromPayPal
// 接收paypal订单事件-webhook接口
// /api/v1/player//tarot/webhook/paypal
func ReceiveOrderEventsFromPayPal(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logging.Debugf("read Webhook error: %v", err)
		return
	}
	var paypalWk WebHookOfPayPalApprovedStruct
	err = json.Unmarshal(body, &paypalWk)
	if err != nil {
		logging.Debugf("Unmarshal Webhook error: %v", err)
		return
	}
	logging.Debugf("webhook detail: %v", paypalWk)
	logging.Debugf("paypalWk.Resource: %v", paypalWk.Resource)
	logging.Debugf("paypalWk.Resource.Id: %v", paypalWk.Resource.Id)
	logging.Debugf("paypalWk.Resource.Intent : %v", paypalWk.Resource.Intent)
	logging.Debugf("paypalWk.Resource.Status: %v", paypalWk.Resource.Status)
	logging.Debugf("paypalWk.Resource.Links: %v", paypalWk.Resource.Links)

	appG := app.Gin{C: c}
	orderId := paypalWk.Resource.Id

	err = CaptureOrder(orderId)
	if err != nil {
		appG.Response(http.StatusOK, "确认订单失败", nil)
		return
	}
	logging.Debugf("OrderID :%s 支付成功!", orderId)
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func CaptureOrder(orderId string) (err error) {
	client, err := paypal.NewClient(PayPalClientID, PayPalSecret, false)
	if err != nil {
		xlog.Error(err)
		logging.Debugf("Error %v", err)
		return
	}
	// 打开Debug开关，输出日志
	client.DebugSwitch = gopay.DebugOn
	ctx := context.Background()
	ppRspc, err := client.OrderCapture(ctx, orderId, nil)
	if err != nil {
		xlog.Error(err)
		logging.Debugf("Error %v", err)
		return
	}
	if ppRspc.Code != paypal.Success {
		// TODO ？？
		logging.Debugf("Retuen Code %v", ppRspc.Code)
		return
	}

	// TODO：判断订单状态 COMPLETED 代表完成

	transaction := ppRspc.Response.PurchaseUnits[0].Payments.Captures[0]
	if transaction.Status != "COMPLETED" {
		logging.Debugf("交易单号：%s", transaction.Id)
		return errors.New(fmt.Sprintf("Status is: %s", transaction.Status))
	}
	// 记录订单交易单号
	// TODO 更新订单状态为已支付，并在订单表中记录交易单号

	err = order_service.UpdateOrderStatus(orderId, transaction.Status, transaction.Id)
	if err != nil {
		logging.Debugf("更新订单状态失败：%v", err)
		return
	}
	return nil
}

type WebHookOfPayPalApprovedStruct struct {
	ID              string            `json:"id"`
	CreateTime      string            `json:"create_time"`
	ResourceType    string            `json:"resource_type"`
	EventType       string            `json:"event_type"`
	Summary         string            `json:"summary"`
	Resource        ResourceOfWebhook `json:"resource"`
	Links           []interface{}     `json:"links"`
	EventVersion    string            `json:"event_version"`
	ResourceVersion string            `json:"resource_version"`
}

type WebHookOfPayPalCompletedStruct struct {
	ID              string            `json:"id"`
	CreateTime      string            `json:"create_time"`
	ResourceType    string            `json:"resource_type"`
	EventType       string            `json:"event_type"`
	Summary         string            `json:"summary"`
	Resource        ResourceOfWebhook `json:"resource"`
	Links           []interface{}     `json:"links"`
	EventVersion    string            `json:"event_version"`
	Zts             int               `json:"zts"`
	ResourceVersion string            `json:"resource_version"`
}

type ResourceOfWebhook struct {
	UpdateTime    string                  `json:"update_time"`
	CreateTime    string                  `json:"create_time"`
	PurchaseUnits []PurchaseUnitOfWebhook `json:"purchase_units"`
	Links         []paypal.Link           `json:"links"`
	Id            string                  `json:"id"`
	GrossAmount   paypal.Amount           `json:"gross_amount"`
	Intent        string                  `json:"intent"`
	Payer         paypal.Payer            `json:"payer"`
	Status        string                  `json:"status"`
}

type PurchaseUnitOfWebhook struct {
	ReferenceId string             `json:"reference_id,omitempty"`
	Amount      *paypal.Amount     `json:"amount,omitempty"`
	Payee       *PayeeOfWebhook    `json:"payee,omitempty"`
	Shipping    *ShippingOfWebhook `json:"shipping,omitempty"`
	//Payments    *PaymentsOfWebhook `json:"payments,omitempty"`
}

type PayeeOfWebhook struct {
	EmailAddress string `json:"email_address,omitempty"`
}

type ShippingOfWebhook struct {
	Method  string `json:"method"`
	Address *AddressOfWebhook
}

type AddressOfWebhook struct {
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	AdminArea1   string `json:"admin_area_1"`
	AdminArea2   string `json:"admin_area_2"`
	PostalCode   string `json:"postal_code"`
	CountryCode  string `json:"country_code"`
}

type PaymentsOfWebhook struct {
	Captures []*paypal.Capture `json:"captures,omitempty"`
}
