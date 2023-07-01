package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/pkg/mail"
)

type SendMailForm struct {
	SendTo    string `json:"send_to" binding:"required"`              // 接受者
	SendType  string `json:"send_type" binding:"required" enums:"ta"` // 发送类型，ta:塔罗订单
	OrderId   string `json:"order_id" binding:"required"`             // 订单号
	AnswerUrl string `json:"answer_url"`                              // 塔罗解答URL
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
	//ansDate := time.Unix(order.PayedTime, 0).Format("2006-01-02")
	// ta 代表塔罗
	if form.SendType == "ta" {

		tarots, _, err := models.GetTarots(0, 10000, "id = ? ", []interface{}{order.TarotList[0]})
		language := tarots[0].Language
		subject := ""
		HTML := ""
		if language == "zh" || language == "tc" {
			if form.AnswerUrl == "" {
				form.AnswerUrl = "https://www.kiminouchuu.com/zh-cn/tarot-reading"
			}
			order.PayMethod = "苹果支付"
			subject = fmt.Sprintf("【小小的宇宙】%s塔罗占卜解答编号!", questDate)
			HTML = fmt.Sprintf("<!DOCTYPE html>\n<html>\n<head> <meta charset=\"utf-8\">\n</head>\n\n小宇宙，大感觉。 当下最受欢迎的算命网页! 在这里--找出关于这个世界、宇宙和你的命运。 从今天的运势到你和他的爱情生活。\n<br><br>\n\n你好，小宇宙。\n<br><br>\n\n谢谢你对塔罗牌答案的兴趣。以下是订购信息。\n<br><br>\n<h3>▶ 送信情報</h3>\n占卜日期：%s\n订单编号：%s\n<br>\n\n<h3>▶ 支付信息</h3>\n支付方式：%s\n电子邮件：%s\n<br><br>\n\n\n<a href=\"%s\">小宇宙 - 塔罗牌答案</a>页面，在那里你可以检查塔罗牌答案的提交状态，并了解你可能有的问题。\n<br>\n\n<h3>▶ 关于我们・小小的宇宙</h3>\n2023年塔罗牌算命网站在最受欢迎的时候，讲述今天的运势。\n爱情运势和塔罗牌\n每天抽一次签，获得幸运。\n占卜和算命。成为你生活中的领导者! 获得幸运!\n至今已有超过20,000人。\n<br><br>\n\n这一次，许个愿吧!\n<br><br><br>\n\nCopyright ⓒ 2022 - 2023 All rights reserved. kiminouchuu.com\n\n</html>", questDate, order.OrderId, order.PayMethod, form.SendTo, form.AnswerUrl)
		} else if language == "en" {
			if form.AnswerUrl == "" {
				form.AnswerUrl = "https://www.kiminouchuu.com/en/tarot-reading"
			}
			subject = fmt.Sprintf("【Tiniverse】%sTarot divination answer number!!", questDate)
			HTML = fmt.Sprintf("<!DOCTYPE html>\n<html>\n<head> <meta charset=\"utf-8\">\n</head>\n\nSmall universe, big feeling. The most popular fortune telling page of the moment! Here - find out about the world, the universe and your destiny. From today's horoscope to your love life with him.\n<br><br>\n\nHello, tiniverse\n<br><br>\n\nThank you for your interest in Tarot answers. Here is the ordering information.\n<br><br>\n<h3>▶ 送信情報</h3>\nDate of divination：%s\nOrder Number：%s\n<br>\n\n<h3>▶ Payment Info</h3>\nPayMethod：%s\nEmail：%s\n<br><br>\n\n\n<a href=\"%s\">Tiniverse - Tarot Answers Page</a>There you can check the submission status of your tarot answers and learn about any questions you may have.\n<br>\n\n<h3>▶ About Us・Tiniverse</h3>\nThe 2023 Tarot fortune telling site tells today's horoscope at the most popular time of the year.\nLove horoscopes and tarot cards\nDraw once a day to get lucky.\nDivination and fortune telling. Be the leader in your life! Get lucky!\nOver 20,000 people so far.\n<br><br>\n\nThis time, make a wish!\n<br><br><br>\n\nCopyright ⓒ 2022 - 2023 All rights reserved. kiminouchuu.com\n\n</html>", questDate, order.OrderId, order.PayMethod, form.SendTo, form.AnswerUrl)
		} else {
			if form.AnswerUrl == "" {
				form.AnswerUrl = "https://www.kiminouchuu.com/tarot-reading"
			}
			subject = fmt.Sprintf("【小さな宇宙】%sでタロット回答注文番号!", questDate)
			HTML = fmt.Sprintf("<!DOCTYPE html>\n<html>\n<head> <meta charset=\"utf-8\">\n</head>\n\n小さな宇宙、大きな感じ。今一番流行っている占いページ！\nここでは―この世界、宇宙、そして自分の運命を知る。今日の運勢から彼との恋愛まで。\n<br><br>\n\nこんにちは、小さな宇宙 です。\n<br><br>\n\nこの度はタロットアンサーをご利用いただきありがとうございます。\n以下、注文情報です。\n<br><br>\n<h3>▶ 送信情報</h3>\nタロットの回答提出日：%s\n注文番号：%s\n<br>\n\n<h3>▶ お支払いについて</h3>\n支払方法：%s\n支払いメールです：%s\n<br><br>\n\n<a href=\"%s\">小さな宇宙-タロット解答</a>のページでは、タロット解答の投稿状況を確認することができ、疑問点を把握することができます。\n<br>\n\n<h3>▶ 私たちについて・小さな宇宙</h3>\n2023 年人気絶頂のタロット占いサイト、今日の運勢を占う\n恋愛占いとタロットカード\n毎日一発の抽選で幸運を得ます\n占いや占いが好き。 あなたの人生のリーダーになりましょう！ 幸運を手に入れよう\nこれまで2万人以上の実績\n<br><br>\n\n今回は、願いを叶える！\n<br><br><br>\n\nCopyright ⓒ 2022 - 2023 All rights reserved. kiminouchuu.com\n\n</html>\n", questDate, order.OrderId, order.PayMethod, form.SendTo, form.AnswerUrl)
		}

		err = mail.SendMail(receiverList, subject, HTML)
		if err != nil {
			logging.Error(fmt.Sprintf("发送邮件出错：%s", err.Error()))
			appG.Response(http.StatusOK, "发送失败", nil)
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	}

}
