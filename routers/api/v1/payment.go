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
	PayPalClientID = "AdLh3w0FtvsXL9LLDfgKY5jJUt0C8LxEjpU0tvKMxIOFJ08tAtp9H4680zb3P2hxxrMaV-iX8CNVojzq"
	PayPalSecret   = "EDg0Z3QrZVyYsqapxoaY9oh9pXjs98qDlr0By23h9vxRoBmN0NTX4UzafXgCXRX6yDDb7BBFtf4K0APo"
)

type PayPalOrderResp struct {
	OrderID string `json:"order_id"`
}

type CreatePayPalOrderForm struct {
	CurrencyCode string `json:"currency_code" binding:"required"` // "货币代码"
	Amount       int    `json:"amount" binding:"required"`        // "金额"
	Uid          string `json:"uid" binding:"required"`           // "用户ID"
}

// CreatePayPalOrder
// @Summary 创建PapPal支付订单
// @Accept json
// @Produce  json
// @Param _ body CreatePayPalOrderForm true "创建订单请求参数"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /player/paypal/checkout/order [post]
// @Tags Player
func CreatePayPalOrder(c *gin.Context) {
	appG := app.Gin{C: c}
	var formD CreatePayPalOrderForm
	if err := c.ShouldBindJSON(formD); err != nil {
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
	// Create Orders example
	var pus []*paypal.PurchaseUnit
	var item = &paypal.PurchaseUnit{
		ReferenceId: "TX12333331231232",
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
				Set("locale", "en-PT").
				Set("return_url", "https://example.com/returnUrl").
				Set("cancel_url", "https://example.com/cancelUrl")
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
	PayPalOrderId := "TA" + ppRsp.Response.Id + ts
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
// @Router /player/paypal/capture/order/{order_id} [post]
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
