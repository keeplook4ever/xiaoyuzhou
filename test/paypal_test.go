package test

import (
	"context"
	"github.com/go-pay/gopay/paypal"
	"testing"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/xlog"
)

func TestGoPay(t *testing.T) {
	xlog.Info("GoPay Version: ", gopay.Version)
}

var (
	Clientid = "AdLh3w0FtvsXL9LLDfgKY5jJUt0C8LxEjpU0tvKMxIOFJ08tAtp9H4680zb3P2hxxrMaV-iX8CNVojzq"
	Secret   = "EDg0Z3QrZVyYsqapxoaY9oh9pXjs98qDlr0By23h9vxRoBmN0NTX4UzafXgCXRX6yDDb7BBFtf4K0APo"
)

//func TestMain(m *testing.M) {
//
//	os.Exit(m.Run())
//}

func TestCreateOrder(t *testing.T) {

	client, err := paypal.NewClient(Clientid, Secret, false)
	if err != nil {
		t.Error(err)
		return
	}
	// 打开Debug开关，输出日志
	client.DebugSwitch = gopay.DebugOn

	xlog.Debugf("Appid: %s", client.Appid)
	xlog.Debugf("AccessToken: %s", client.AccessToken)
	xlog.Debugf("ExpiresIn: %d", client.ExpiresIn)

	// Create Orders example
	var pus []*paypal.PurchaseUnit
	var item = &paypal.PurchaseUnit{
		ReferenceId: "TX12333331231232",
		Amount: &paypal.Amount{
			CurrencyCode: "USD",
			Value:        "8",
		},
	}
	pus = append(pus, item)

	bm := make(gopay.BodyMap)
	bm.Set("intent", "CAPTURE").
		Set("purchase_units", pus).
		SetBodyMap("application_context", func(b gopay.BodyMap) {
			b.Set("brand_name", "gopay").
				Set("locale", "en-PT").
				Set("return_url", "https://example.com/returnUrl").
				Set("cancel_url", "https://example.com/cancelUrl")
		})
	ctx := context.Background()
	ppRsp, err := client.CreateOrder(ctx, bm)
	if err != nil {
		xlog.Error(err)
		return
	}
	if ppRsp.Code != paypal.Success {
		// do something
		xlog.Error("!!!")
		return
	}

	ppRspc, err := client.OrderCapture(ctx, ppRsp.Response.Id, nil)
	if err != nil {
		xlog.Error(err)
		return
	}
	if ppRspc.Code != paypal.Success {
		// do something
		return
	}
}
