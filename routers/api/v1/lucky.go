package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/setting"
	"xiaoyuzhou/pkg/util"
	"xiaoyuzhou/service/lucky_service"
)

type AddLuckyForm struct {
	Spell string `json:"spell" binding:"required"`
	Todo  string `json:"todo" binding:"required"`
	Song  string `json:"song" binding:"required"`
}

type EditLuckyForm struct {
	Id    int    `json:"id" binding:"required"`
	Spell string `json:"spell"`
	Todo  string `json:"todo"`
	Song  string `json:"song"`
}

// AddLucky
// @Summary 添加今日好运内容
// @Param _ body AddLuckyForm true "参数"
// @Produce json
// @Accept json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lucky [post]
// @Security ApiKeyAuth
// @Tags Manager
func AddLucky(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqData AddLuckyForm
	if err := c.ShouldBindJSON(&reqData); err != nil {
		appG.Response(http.StatusBadRequest, "非法请求", nil)
		return
	}
	luckyI := lucky_service.LuckyInput{
		Spell: reqData.Spell,
		Todo:  reqData.Todo,
		Song:  reqData.Song,
	}
	if err := luckyI.Add(); err != nil {
		appG.Response(http.StatusOK, err.Error(), nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// EditLucky
// @Summary 修改今日好运内容
// @Param _ body EditLuckyForm true "参数"
// @Produce json
// @Accept json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lucky/{id} [put]
// @Security ApiKeyAuth
// @Tags Manager
func EditLucky(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqData EditLuckyForm
	if err := c.ShouldBindJSON(&reqData); err != nil {
		appG.Response(http.StatusBadRequest, "参数不合法", nil)
		return
	}
	luckI := lucky_service.LuckyInput{
		Id:    com.StrTo(c.Param("id")).MustInt(),
		Spell: reqData.Spell,
		Todo:  reqData.Todo,
		Song:  reqData.Song,
	}
	if err := luckI.Edit(); err != nil {
		appG.Response(http.StatusOK, err.Error(), nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// DeleteLucky
// @Summary 修改今日好运内容
// @Param id path int true "ID"
// @Produce json
// @Accept json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lucky/{id} [delete]
// @Security ApiKeyAuth
// @Tags Manager
func DeleteLucky(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	luckI := lucky_service.LuckyInput{
		Id: id,
	}
	if err := luckI.Delete(); err != nil {
		appG.Response(http.StatusOK, err.Error(), nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// GetLucky
// @Summary 获取今日好运内容
// @Param spell query string false "咒语"
// @Param todo query string false "适宜"
// @Param song query string false "歌曲"
// @Produce json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lucky [get]
// @Security ApiKeyAuth
// @Tags Manager
func GetLucky(c *gin.Context) {
	appG := app.Gin{C: c}
	luckI := lucky_service.LuckyInput{
		Spell:    c.Query("spell"),
		Todo:     c.Query("todo"),
		Song:     c.Query("song"),
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	luckys, err := luckI.Get()
	if err != nil {
		appG.Response(http.StatusOK, err.Error(), nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, GetLuckyResponse{Lists: luckys, Count: len(luckys)})
}

type GetLuckyResponse struct {
	Lists []models.LuckyTodayDto `json:"lists"`
	Count int                    `json:"count"`
}
