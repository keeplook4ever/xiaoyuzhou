package player

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
)

type TarotFree struct {
	Id     int
	ImgUrl string
	Name   string
	Type_  string
}

type TarotPaid struct {
	BasicT   TarotFree
	CardRead string
	TypeRead string
	Answer   string
}

// GetTarotFree
// @Summary 免费获取塔罗牌
// @Produce  json
// @Param uid query int true "用户id"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /player/tarot [get]
// @Tags Player
func GetTarotFree(c *gin.Context) {
	appG := app.Gin{C: c}
	tarot := TarotFree{
		Id:     1,
		ImgUrl: "",
		Name:   "",
		Type_:  "",
	}
	appG.Response(http.StatusOK, e.SUCCESS, tarot)
}

// GetTarotPaid
// @Summary 获取塔罗牌付费内容
// @Produce  json
// @Param uid body int true "用户id"
// @Param tarotId body string true "塔罗牌id"
// @Param paidOrderId body string true "付费id"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /player/tarot/read [post]
// @Tags Player
func GetTarotPaid(c *gin.Context) {

	// TODO: 需要验证支付单是否已经支付成功，有效
	// TODO: 请求参数验证

	appG := app.Gin{C: c}
	var tarot TarotPaid
	tarot.BasicT.Id = 1
	tarot.BasicT.Name = ""
	tarot.BasicT.ImgUrl = ""
	tarot.BasicT.Type_ = ""
	tarot.TypeRead = ""
	tarot.CardRead = ""
	tarot.Answer = ""
	appG.Response(http.StatusOK, e.SUCCESS, tarot)
}
