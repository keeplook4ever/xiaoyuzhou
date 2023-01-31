package payment_service

import (
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/xlog"
)

func TestGoPay() {
	xlog.Info("GoPay Version: ", gopay.Version)
}
