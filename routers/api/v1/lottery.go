package v1

import (
	"github.com/gin-gonic/gin"
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

type EditLotteryContentForm struct {
	Content string `form:"content" binding:"required"`
	Type      string    `form:"type" binding:"required"`
	Id int `path:"id" binding:"required"`
}

type EditLotteryForm struct {
	MaxScore int `form:"max_score"`
	MinScore int `form:"min_score"`
	Probability float32 `form:"probability"`
	KeyWord string `form:"keyword"`
	Type string `form:"type" binding:"required" enums:"A,B,C,D"`// A-D 枚举
}

type AddLotteryContentData struct {
	Type string `json:"type" binding:"required" enums:"A,B,C,D"` //枚举A-D
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


// EditLottery
// @Summary 修改运势类型
// @Produce json
// @Param max_score formData int false "最大值"
// @Param min_score formData int false "最小值"
// @Param probability formData float32 false "概率"
// @Param keyword formData string false "关键字"
// @Param type formData string true "好运等级" Enums("A","B","C","D")
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lottery [put]
// @Tags Manager
// @Security ApiKeyAuth
func EditLottery(c *gin.Context) {
	appG := app.Gin{C: c}
	var l EditLotteryForm
	if err := c.ShouldBind(&l); err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	lotteryInput := lottery_service.LotteryInput{
		MaxScore: l.MaxScore,
		MinScore: l.MinScore,
		KeyWord: l.KeyWord,
		Type: l.Type,
		Probability: l.Probability,
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
		Content: lotteryC.Content,
		Type: lotteryC.Type,  // type 代表A-D不同等级
	}
	if err := lotteryContentInput.Add(); err != nil {
		appG.Response(http.StatusOK, "添加LotteryContent失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// EditLotteryContent
// @Summary 修改运势详细内容
// @Produce json
// @Param id path int true "ID"
// @Param type formData string true "好运等级" Enums("A","B","C","D")
// @Param content formData string true "Content"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lottery-content/{id} [put]
// @Tags Manager
// @Security ApiKeyAuth
func EditLotteryContent(c *gin.Context) {
	appG := app.Gin{C: c}
	var Lc = EditLotteryContentForm{
		//Id: com.StrTo(c.Param("id")).MustInt(),
	}
	if err := c.ShouldBind(&Lc); err != nil {
		appG.Response(http.StatusBadRequest, "请求不合法", nil)
		return
	}

	lcInput := lottery_service.LotteryContentInput{
		ID: 		Lc.Id,
		Type:      	Lc.Type,
		Content: 	Lc.Content,
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
// @Param type query string false "好运等级" Enums("A","B","C","D")
// @Success 200 {object} []models.LotteryContent
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lottery-content [get]
// @Tags Manager
// @Security ApiKeyAuth
func GetLotteryContentForManager(c *gin.Context) {
	appG := app.Gin{C: c}
	tP := c.Query("type")
	lotteryInput := lottery_service.LotteryContentInput{
		Type: tP,
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
