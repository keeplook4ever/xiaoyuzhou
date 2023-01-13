package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"reflect"
	"sort"
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

type InputLotteryData struct {
	KeyWordList     []string  `json:"keyword_list" binding:"required"`     //["末吉", "小吉", "吉", "大吉"]
	ScoreList       []int     `json:"score_list" binding:"required"`       //从小到大
	ProbabilityList []float32 `json:"probability_list" binding:"required"` //相加为1
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
// @Success 200 {object} []models.Lottery
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lottery [get]
// @Tags Manager
// @Security ApiKeyAuth
func GetLotteryForManager(c *gin.Context) {
	appG := app.Gin{C: c}
	var lotteries []models.Lottery
	lotteries, err := lottery_service.GetLotteryForManager()
	if err != nil {
		appG.Response(http.StatusOK, "获取运势表出错", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, lotteries)
}

// AddLotteryType
// @Summary 添加运势类型关键字，分数，概率等
// @Produce json
// @Accept json
// @Param _ body InputLotteryData true "运势类"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lottery [post]
// @Tags Manager
// @Security ApiKeyAuth
func AddLotteryType(c *gin.Context) {
	appG := app.Gin{C: c}
	var lottery InputLotteryData
	if err := c.ShouldBindJSON(&lottery); err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	//校验上传值是否合法
	if !checkLotteryValid(lottery.KeyWordList, lottery.ScoreList, lottery.ProbabilityList) {
		appG.Response(http.StatusOK, "参数不合法", nil)
		return
	}
	lotteryInput := lottery_service.LotteryInput{
		KeyWordList:     lottery.KeyWordList,
		ScoreList:       lottery.ScoreList,
		ProbabilityList: lottery.ProbabilityList,
	}
	if err := lotteryInput.Add(); err != nil {
		appG.Response(http.StatusOK, "添加Lottery出错", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// EditLottery
// @Summary 修改运势类型
// @Produce json
// @Param _ body InputLotteryData true "运势类"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lottery [put]
// @Tags Manager
// @Security ApiKeyAuth
func EditLottery(c *gin.Context) {
	appG := app.Gin{C: c}
	var l InputLotteryData
	if err := c.ShouldBind(&l); err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	lotteryInput := lottery_service.LotteryInput{
		KeyWordList:     l.KeyWordList,
		ScoreList:       l.ScoreList,
		ProbabilityList: l.ProbabilityList,
	}

	if err := lotteryInput.Edit(); err != nil {
		appG.Response(http.StatusOK, "编辑失败", nil)
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

// EditLotteryContent
// @Summary 添加运势详细内容
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param keyword formData string false "KeyWord"
// @Param content formData string false "Content"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lottery-content/[id] [put]
// @Tags Manager
// @Security ApiKeyAuth
func EditLotteryContent(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	keyWord := c.PostForm("keyword")
	content := c.PostForm("content")
	lcInput := lottery_service.LotteryContentInput{
		ID:      id,
		KeyWord: keyWord,
		Content: content,
	}
	if err := lcInput.Update(); err != nil {
		appG.Response(http.StatusOK, "更新运势内容失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

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
	lotteryInput := lottery_service.LotteryContentInput{
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

// 校验上传Lottery是否合法
func checkLotteryValid(keywordList []string, scoreList []int, probList []float32) bool {
	// 校验数量是否一致
	if len(keywordList) != len(scoreList) || len(keywordList) != len(probList) || len(scoreList) != len(probList) {
		return false
	}

	// 校验score是否从小到大
	sortedScoreList := make([]int, len(scoreList))
	copy(sortedScoreList, scoreList)
	sort.Ints(sortedScoreList)
	if !reflect.DeepEqual(scoreList, sortedScoreList) {
		return false
	}

	// 校验probList 概率相加是否等于1
	totalValue := float32(0)
	for _, v := range probList {
		totalValue = v + totalValue
	}
	if totalValue != 1 {
		return false
	}
	return true
}
