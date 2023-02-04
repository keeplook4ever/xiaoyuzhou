package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/paypal"
	"github.com/go-pay/gopay/pkg/xlog"
	"net/http"
	"strconv"
	"time"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
)

const (
	PayPalClientID = "AYxax5vAtOVCwWE7k3-IkfeeW91WssDy0X87o3hYU6v64yJPWDkb_9mWbcKrlixMDssEXnaU75Qwd8A3"
	PayPalSecret   = "EIIDNX57BlG8pCeFiiY6WJipeMtSQsIiOOP3ojg0_gtNSd3ndB0asdBQXP9IIeGh6gmlRnJHjjeozKGP"
)

type PayPalOrderResp struct {
	OrderID string `json:"order_id"`
}

type CreatePayPalOrderForm struct {
	CurrencyCode string `json:"currency_code" binding:"required"` // "货币代码"
	Amount       int    `json:"amount" binding:"required"`        // "金额"
	Uid          string `json:"uid" binding:"required"`           // "用户ID"
	ReturnURL    string `json:"return_url" binding:"required"`    // "支付成功返回URL"
	CancelURL    string `json:"cancel_url"  binding:"required"`   // "取消支付URL"
}

// 账单地址：

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
	var item = &paypal.PurchaseUnit{
		//ReferenceId: "TX12333331231232",
		//TODO:
		Amount: &paypal.Amount{
			CurrencyCode: formD.CurrencyCode,
			Value:        strconv.Itoa(formD.Amount),
		},
	}
	pus = append(pus, item)

	bm := make(gopay.BodyMap)
	bm.Set("intent", "CAPTURE").
		Set("purchase_units", pus).
		SetBodyMap("application_context", func(b gopay.BodyMap) {
			b.Set("brand_name", "xiaoyuzhou").
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

	//TODO：将订单落库

	ts := strconv.Itoa(int(time.Now().Unix()))
	PayPalOrderId := "TA" + "-" + ppRsp.Response.Id + "-" + ts
	//TODO: 存入数据库
	xlog.Debugf("OrderID: %s", PayPalOrderId)

	appG.Response(http.StatusOK, e.SUCCESS, ppRsp.Response)
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
func CapturePayPalOrder(c *gin.Context) {
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
	ppRspc, err := client.OrderCapture(ctx, orderId, nil)
	if err != nil {
		xlog.Error(err)
		appG.Response(http.StatusOK, "捕获订单失败", nil)
		return
	}
	if ppRspc.Code != paypal.Success {
		// TODO ？？

		appG.Response(http.StatusOK, "捕获订单失败", nil)
		return
	}

	// TODO：判断订单状态 COMPLETED 代表完成
	appG.Response(http.StatusOK, e.SUCCESS, ppRspc.Response)
}

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
