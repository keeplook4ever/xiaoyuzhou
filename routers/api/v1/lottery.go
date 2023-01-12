package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/service/lottery_service"
)

type GetLotteryResponse struct {
	LotteryContent models.LotteryDto
	LuckyContent   models.LuckyTodayDto
}

type AddLotteryData struct {
	MaxScore    int     `json:"max_score" binding:"required"`
	MinScore    int     `json:"min_score" bingding:"required"`
	KeyWord     string  `json:"key_word" binding:"required"`
	Probability float32 `json:"probability" binding:"required"`
}

type AddLotteryContentData struct {
	KeyWord string `json:"key_word" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// GetLottery
// @Summary 获取日签
// @Produce  json
// @Param uid query string true "用户uid"
// @Success 200 {object} GetLotteryResponse
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /player/lottery [get]
// @Tags Player

func GetLottery(c *gin.Context) {
	appG := app.Gin{C: c}
	uid := c.Query("uid")
	logging.Debugf("uid is: %s", uid)
	// 存储用户uid的记录
	lottery, err := lottery_service.GetLotteryForPlayer()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorGetLotteryFail, nil)
		return
	}
	luckyTody, err := lottery_service.GetLuckyForPlayer()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorGetLuckytodyFail, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, GetLotteryResponse{LotteryContent: lottery, LuckyContent: luckyTody})
}

// GetLotteryForManager
// @Summary 获取全部运势表Lottery
// @Produce json
// @Param keyword formData string false "运势文字"
// @Success 200 {object} []models.Lottery
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lottery [get]
// @Tags Manager
// @Security ApiKeyAuth
func GetLotteryForManager(c *gin.Context) {
	appG := app.Gin{C: c}
	keyword := c.Query("keyword")
	lotteryInput := lottery_service.LotteryInput{
		KeyWord: keyword,
	}
	var lotteries []models.Lottery
	lotteries, err := lotteryInput.GetLotteryForManager()
	if err != nil {
		appG.Response(http.StatusOK, "获取运势表出错", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, lotteries)
}

// GetLotteryContentForManager
// @Summary 获取全部运势内容表LotteryContent
// @Produce json
// @Param keyword formData string false "运势文字"
// @Success 200 {object} []models.LotteryContent
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lottery-content [get]
// @Tags Manager
// @Security ApiKeyAuth
func GetLotteryContentForManager(c *gin.Context) {
	appG := app.Gin{C: c}
	keyword := c.Query("keyword")
	lotteryInput := lottery_service.LotteryInput{
		KeyWord: keyword,
	}

	var lotteryContents []models.LotteryContent
	lotteryContents, err := lotteryInput.GetLotteryContentForManager()
	if err != nil {
		appG.Response(http.StatusOK, "获取运势内容表出错", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, lotteryContents)
}

// AddLotteryType
// @Summary 添加运势类型关键字，分数，概率等
// @Produce json
// @Accept json
// @Param _ body AddLotteryData true "运势类关键"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lottery [post]
// @Tags Manager
// @Security ApiKeyAuth
func AddLotteryType(c *gin.Context) {
	appG := app.Gin{C: c}
	var lottery AddLotteryData
	if err := c.ShouldBindJSON(&lottery); err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	lotteryInput := lottery_service.LotteryInput{
		MaxScore:    lottery.MaxScore,
		MinScore:    lottery.MinScore,
		KeyWord:     lottery.KeyWord,
		Probability: lottery.Probability,
	}
	if err := lotteryInput.Add(); err != nil {
		appG.Response(http.StatusOK, "添加Lottery出错", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// AddLotteryContent
// @Summary 添加运势详细内容
// @Accept json
// @Produce json
// @Param _ body AddLotteryContentData true "运势具体内容"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lottery-content [post]
// @Tags Manager
// @Security ApiKeyAuth
func AddLotteryContent(c *gin.Context) {
	appG := app.Gin{C: c}
	var lotteryC AddLotteryContentData
	if err := c.ShouldBindJSON(&lotteryC); err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	lotteryContentInput := lottery_service.LotteryContentInput{
		KeyWord: lotteryC.KeyWord,
		Content: lotteryC.Content,
	}
	if err := lotteryContentInput.Add(); err != nil {
		appG.Response(http.StatusOK, "添加LotteryContent失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditLotteryForm struct {
	MaxScore    int     `form:"max_score"`
	MinScore    int     `form:"min_score"`
	KeyWord     string  `form:"keyword"`
	Probability float32 `form:"probability"`
	ID          int     `form:"id" binding:"required"`
}

// EditLottery
// @Summary 修改运势类型
// @Produce json
// @Param max_score formData int false "最大分数"
// @Param min_score formData int false "最小分数"
// @Param keyword formData string false "关键字"
// @Param probability formData float32 false "概率"
func EditLottery(c *gin.Context) {
	appG := app.Gin{C: c}
	var l EditLotteryForm
	if err := c.ShouldBindJSON(&l); err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	lotteryInput := lottery_service.LotteryInput{
		KeyWord:     l.KeyWord,
		MaxScore:    l.MaxScore,
		MinScore:    l.MinScore,
		Probability: l.Probability,
		ID:          l.ID,
	}

	if err := lotteryInput.Edit(); err != nil {
		appG.Response(http.StatusOK, "编辑失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
