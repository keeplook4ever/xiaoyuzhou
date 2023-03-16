package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/mail"
)

type SendMailForm struct {
	SendTo   string `json:"send_to" binding:"required"`              // 接受者
	SendType string `json:"send_type" binding:"required" enums:"ta"` // 发送类型，ta:塔罗订单
	OrderId  string `json:"order_id" binding:"required"`             // 订单号
	Uid      string `json:"uid" binding:"required"`                  // 用户id
}

// SendMailToCustomer
// @Summary 发送邮件给客户
// @Param _ body SendMailForm true "发送邮件参数"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /player/mail [post]
// @Tags Player
func SendMailToCustomer(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form SendMailForm
	)
	if err := c.ShouldBindJSON(&form); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}
	receiverList := []string{
		form.SendTo,
	}

	//content, err := order_service.GetOrderForMail(form.Uid, form.OrderId)

	// 如果是塔罗牌发邮件
	if form.SendType == "ta" {
		subject := "Here is your Tarot Order"
		txt := "Welcome to mailgun and kiminouchuu.com"
		err := mail.SendMail(receiverList, subject, txt)
		if err != nil {
			appG.Response(http.StatusOK, "发送失败", nil)
			return
		}
	}

}
