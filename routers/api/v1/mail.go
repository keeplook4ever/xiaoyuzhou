package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/mail"
)

type SendMailForm struct {
	SendTo    string `json:"send_to" binding:"required"`              // 接受者
	SendType  string `json:"send_type" binding:"required" enums:"ta"` // 发送类型，ta:塔罗订单
	OrderId   string `json:"order_id" binding:"required"`             // 订单号
	Uid       string `json:"uid" binding:"required"`                  // 用户id
	AnswerUrl string `json:"answer_url" binding:"required"`           // 塔罗解答URL
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

	order, err := models.GetOrderByOrderId(form.OrderId)
	if err != nil {
		appG.Response(http.StatusOK, "获取订单失败", nil)
		return
	}
	questDate := time.Unix(order.PickTime, 0).Format("2006-01-02")
	ansDate := time.Unix(order.PayedTime, 0).Format("2006-01-02")
	// ta 代表塔罗
	if form.SendType == "ta" {
		subject := fmt.Sprintf("【小さな宇宙】%sでタロット回答注文番号!", questDate)
		JPHTML := fmt.Sprintf("<!DOCTYPE html>\n<html>\n<head> <meta charset=\"utf-8\">\n</head>\n\n小さな宇宙、大きな感じ。今一番流行っている占いページ！\nここでは―この世界、宇宙、そして自分の運命を知る。今日の運勢から彼との恋愛まで。\n<br><br>\n\nこんにちは、小さな宇宙 です。\n<br><br>\n\nこの度はタロットアンサーをご利用いただきありがとうございます。\n以下、注文情報です。\n<br><br>\n<h3>▶ 送信情報</h3>\nタロットの回答提出日：%s\n注文番号：%s\n<br>\n\n<h3>▶ お支払いについて</h3>\n支払方法：%s\n支払いメールです：%s\n<br><br>\n\n<a href=\"%s\">小さな宇宙-タロット解答</a>のページでは、タロット解答の投稿状況を確認することができ、疑問点を把握することができます。\n<br>\n\n<h3>▶ 私たちについて・小さな宇宙</h3>\n2023 年人気絶頂のタロット占いサイト、今日の運勢を占う\n恋愛占いとタロットカード\n毎日一発の抽選で幸運を得ます\n占いや占いが好き。 あなたの人生のリーダーになりましょう！ 幸運を手に入れよう\nこれまで2万人以上の実績\n<br><br>\n\n今回は、願いを叶える！\n<br><br><br>\n\nCopyright ⓒ 2022 - 2023 All rights reserved. kiminouchuu.com\n\n</html>\n", ansDate, order.OrderId, order.PayMethod, form.SendTo, form.AnswerUrl)
		err := mail.SendMail(receiverList, subject, JPHTML)
		if err != nil {
			appG.Response(http.StatusOK, "发送失败", nil)
			return
		}
	}

}
