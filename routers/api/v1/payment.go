package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/paypal"
	"io/ioutil"
	"net/http"
	"strconv"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/pkg/setting"
	"xiaoyuzhou/service/order_service"
	"xiaoyuzhou/service/tarot_service"
)

type PayPalOrderResp struct {
	OrderID string `json:"order_id"`
}

type CreatePayPalOrderForm struct {
	OrderId   string `json:"order_id" binding:"required"`    // 订单ID
	ReturnURL string `json:"return_url" binding:"required"`  // 支付成功返回URL
	CancelURL string `json:"cancel_url"  binding:"required"` // 取消支付URL
	// 场景：
	Scene    string `json:"scene" binding:"required" enums:"ta_one_high,ta_one_low,ta_three_high,ta_three_low"` // 支付场景：ta_one_high 塔罗单张高价
	Language string `json:"language" binding:"required" enums:"jp,zh,en,tc"`                                    // 地区:  tc:台湾
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

	// 默认产线模式
	client, err := paypal.NewClient(order_service.PayPalClientIDPrd, order_service.PayPalSecretPrd, true)
	if err != nil {
		logging.Error(fmt.Sprintf("初始化client失败: %s", err.Error()))
		appG.Response(http.StatusOK, "初始化client失败", nil)
		return
	}
	if setting.PaymentSetting.Mode == "debug" {
		client, err = paypal.NewClient(order_service.PayPalClientIDTest, order_service.PayPalSecretTest, false)
		if err != nil {
			logging.Error(fmt.Sprintf("初始化client失败: %s", err.Error()))
			appG.Response(http.StatusOK, "初始化client失败", nil)
			return
		}
	}
	// 打开Debug开关，输出日志
	client.DebugSwitch = gopay.DebugOn

	var pus []*paypal.PurchaseUnit

	// 获取价格
	amount := tarot_service.GetPaymentPrice(formD.Scene, formD.Language)

	// 默认日本
	value := strconv.Itoa(int(amount))
	currencyCode := "JPY"
	localeV := "ja-JP"
	switch formD.Language {
	case "jp":
		currencyCode = "JPY"
		value = strconv.Itoa(int(amount))
		localeV = "ja-JP"
	case "zh":
		currencyCode = "CNY"
		value = fmt.Sprintf("%.2f", amount)
		localeV = "zh-CN"
	case "en":
		currencyCode = "USD"
		value = fmt.Sprintf("%.2f", amount)
		localeV = "en-US"
	case "tc":
		currencyCode = "TWD"
		value = fmt.Sprintf("%.2f", amount)
		localeV = "zh-TW"
	}

	var item = &paypal.PurchaseUnit{
		//ReferenceId: "TX12333331231232",
		Amount: &paypal.Amount{
			CurrencyCode: currencyCode,
			Value:        value,
		},
	}
	pus = append(pus, item)

	bm := make(gopay.BodyMap)
	bm.Set("intent", "CAPTURE").
		Set("purchase_units", pus).
		SetBodyMap("application_context", func(b gopay.BodyMap) {
			b.Set("brand_name", "小さな宇宙").
				Set("locale", localeV).
				Set("return_url", formD.ReturnURL).
				Set("cancel_url", formD.CancelURL)
		})
	ctx := context.Background()
	ppRsp, err := client.CreateOrder(ctx, bm)
	if err != nil {
		logging.Error(fmt.Sprintf("创建订单失败: %s", err.Error()))
		appG.Response(http.StatusOK, "创建订单失败", nil)
		return
	}
	if ppRsp.Code != paypal.Success {
		logging.Error(fmt.Sprintf("创建订单返回码: %v", ppRsp.Code))
		appG.Response(http.StatusOK, "创建订单失败", nil)
		return
	}

	OriOrderId := ppRsp.Response.Id // 原始paypal订单号，后面更新订单状态需要用到
	//将订单落库
	err = order_service.AddPaymentInfoToOrder(formD.OrderId, OriOrderId, "paypal", amount)
	if err != nil {
		logging.Error(fmt.Sprintf("订单落库失败: %s", err.Error()))
		appG.Response(http.StatusOK, "订单落库失败", nil)
		return
	}

	logging.Info(fmt.Sprintf("订单创建成功：%s", formD.OrderId))
	appG.Response(http.StatusOK, e.SUCCESS, ppRsp.Response)
}

// GetOrderStatus
// @Summary 获取订单支付状态
// @Accept json
// @Produce  json
// @Param order_id path string true "订单号"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /player/checkout/orders/{order_id} [get]
// @Tags Player
func GetOrderStatus(c *gin.Context) {
	appG := app.Gin{C: c}
	orderId := c.Param("order_id")
	if orderId == "" {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	//从数据库查订单状态，返回状态: 0,1,2
	status, err := order_service.GetOrderStatus(orderId)
	if err != nil {
		appG.Response(http.StatusOK, "获取订单状态失败,请重试", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, status)
}

// ReceiveOrderEventsFromPayPal
// 接收paypal订单事件-webhook接口 被动查订单，更新订单状态
// /api/v1/player/tarot/webhook/paypal
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
	logging.Debugf("paypalWk.Resource.Id: %v", paypalWk.Resource.Id)
	logging.Debugf("paypalWk.Resource.Intent : %v", paypalWk.Resource.Intent)
	logging.Debugf("paypalWk.Resource.Status: %v", paypalWk.Resource.Status)

	OriOrderId := paypalWk.Resource.Id
	logging.Debugf("OrderID: %s", OriOrderId)
	err = order_service.CaptureOrder(OriOrderId, "paypal")
	if err != nil {
		if err.Error() == "ok" { // ok代表订单状态已确认，或者支付成功，或者失败
			return
		}
		logging.Debugf("Order: %s Failed", OriOrderId)
		logging.Error(fmt.Sprintf("Error Webhook: %s", err.Error()))
		return
	}
	logging.Info(fmt.Sprintf("OrderID :%s 支付成功!", OriOrderId))
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
