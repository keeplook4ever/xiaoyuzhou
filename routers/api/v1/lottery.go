package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/util"
	"xiaoyuzhou/service/lottery_service"
)

type GetLotteryForUserResponse struct {
	LotteryContent models.LotteryDto    `json:"lottery_content"`
	LuckyContent   models.LuckyTodayDto `json:"lucky_content"`
}

type EditLotteryContentForm struct {
	Content  string `form:"content" binding:"required"`
	Type     string `form:"type" binding:"required"`
	Id       int    `path:"id" binding:"required"`
	Language string `form:"language" binding:"required"`
}

type EditLotteryForm struct {
	MaxScore    int     `json:"max_score"`
	MinScore    int     `json:"min_score"`
	Probability float32 `json:"probability"`
	KeyWord     string  `json:"keyword"`
	Type        string  `json:"type" binding:"required" enums:"A,B,C,D"` // A-D 枚举
	//Language    string  `json:"language" enums:"jp,zh,en"`               // 语言
}

type AddLotteryContentData struct {
	Type     string `json:"type" binding:"required" enums:"A,B,C,D"`         //枚举A-D
	Content  string `json:"content" binding:"required"`                      //内容
	Language string `json:"language" binding:"required" enums:"jp,zh,en,tc"` //语言
}

type GetLotteryForManagerResponse struct {
	Lists []models.Lottery `json:"lists"`
	Count int64            `json:"count"`
}

type GetLotteryContentForManagerResponse struct {
	Lists []models.LotteryContent `json:"lists"`
	Count int64                   `json:"count"`
}

// GetLotteryForManager
// @Summary 获取运势表Lottery
// @Produce json
// @Param language query string false "语言" Enums(zh,jp,en,tc)
// @Success 200 {object} GetLotteryForManagerResponse
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lottery [get]
// @Tags Manager
// @Security ApiKeyAuth
func GetLotteryForManager(c *gin.Context) {
	appG := app.Gin{C: c}
	lotteries, count, err := lottery_service.GetLotteryForManager(c.Query("language"))
	if err != nil {
		appG.Response(http.StatusOK, "获取运势表出错", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, GetLotteryForManagerResponse{Lists: lotteries, Count: count})
}

// EditLottery
// @Summary 编辑运势表Lottery
// @Produce json
// @Param _ body []EditLotteryForm true "编辑运势"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lottery [put]
// @Tags Manager
// @Security ApiKeyAuth
func EditLottery(c *gin.Context) {
	appG := app.Gin{C: c}
	var l []EditLotteryForm
	if err := c.ShouldBind(&l); err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}

	if !checkLotteryValid(l) {
		appG.Response(http.StatusOK, "参数校验不通过", nil)
		return
	}

	for _, Lt := range l {
		lotteryI := lottery_service.LotteryInput{
			MaxScore:    Lt.MaxScore,
			MinScore:    Lt.MinScore,
			KeyWord:     Lt.KeyWord,
			Type:        Lt.Type,
			Probability: Lt.Probability,
		}
		if err := lotteryI.Edit(); err != nil {
			appG.Response(http.StatusOK, "编辑失败", nil)
			return
		}
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
		Content:  lotteryC.Content,
		Type:     lotteryC.Type, // type 代表A-D不同等级
		Language: lotteryC.Language,
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
// @Param type formData string true "好运等级" Enums(A,B,C,D)
// @Param content formData string true "Content"
// @Param language formData string true "语言" Enums(zh,jp,en)
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lottery-content/{id} [put]
// @Tags Manager
// @Security ApiKeyAuth
func EditLotteryContent(c *gin.Context) {
	appG := app.Gin{C: c}
	var Lc = EditLotteryContentForm{
		Id: com.StrTo(c.Param("id")).MustInt(),
	}
	if err := c.ShouldBind(&Lc); err != nil {
		appG.Response(http.StatusBadRequest, "请求不合法", nil)
		return
	}

	lcInput := lottery_service.LotteryContentInput{
		ID:       Lc.Id,
		Type:     Lc.Type,
		Content:  Lc.Content,
		Language: Lc.Language,
	}
	if err := lcInput.Update(); err != nil {
		appG.Response(http.StatusOK, "更新运势内容失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// DeleteLotteryContent
// @Summary 删除运势内容
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lottery-content/{id} [delete]
// @Tags Manager
// @Security ApiKeyAuth
func DeleteLotteryContent(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	lcInput := lottery_service.LotteryContentInput{
		ID: id,
	}
	if err := lcInput.Delete(); err != nil {
		appG.Response(http.StatusOK, err.Error(), nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// GetLotteryContentForManager
// @Summary 获取全部运势内容表LotteryContent
// @Produce json
// @Param type query string false "好运等级" Enums(A,B,C,D)
// @Param language query string false "语言" Enums(jp,zh,en,tc)
// @Success 200 {object} GetLotteryContentForManagerResponse
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lottery-content [get]
// @Tags Manager
// @Security ApiKeyAuth
func GetLotteryContentForManager(c *gin.Context) {
	appG := app.Gin{C: c}
	tP := c.Query("type")
	lan := c.Query("language")
	lotteryInput := lottery_service.LotteryContentInput{
		Type:     tP,
		Language: lan,
		PageNum:  util.GetPage(c),
		PageSize: util.GetPageSize(c),
	}

	lotteryContents, count, err := lotteryInput.GetLotteryContentForManager()
	if err != nil {
		appG.Response(http.StatusOK, "获取运势内容表出错", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, GetLotteryContentForManagerResponse{Lists: lotteryContents, Count: count})
}

// 校验上传Lottery是否合法
func checkLotteryValid(editL []EditLotteryForm) bool {
	// 校验数量
	if !(len(editL) == 4) {
		return false
	}

	var (
		AData EditLotteryForm
		BData EditLotteryForm
		CData EditLotteryForm
		DData EditLotteryForm
	)

	for i, v := range editL {
		if v.Type == "A" {
			AData = editL[i]
		} else if v.Type == "B" {
			BData = editL[i]
		} else if v.Type == "C" {
			CData = editL[i]
		} else if v.Type == "D" {
			DData = editL[i]
		}
	}

	// 校验分数 1-99
	if AData.MaxScore != 0 && !(AData.MaxScore < 100 && AData.MinScore == BData.MaxScore+1 && BData.MinScore == CData.MaxScore+1 && CData.MinScore == DData.MaxScore+1 && DData.MinScore > 0) {
		return false
	}

	// 校验概率
	if AData.Probability != 0.0 && !(AData.Probability+BData.Probability+CData.Probability+DData.Probability == 1) {
		return false
	}

	return true
}

// GetLotteryForUser
// @Summary 获取日签
// @Produce  json
// @Param uid query string true "用户uid"
// @Param language query string true "语言" Enums(jp,zh,en,tc)
// @Success 200 {object} GetLotteryForUserResponse
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /player/lottery [get]
// @Tags Player
func GetLotteryForUser(c *gin.Context) {
	appG := app.Gin{C: c}
	uid := c.Query("uid")
	language := c.Query("language")
	lottery, err := lottery_service.GetLotteryForPlayer(uid, language)
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorGetLotteryFail, nil)
		return
	}
	luckyTody, err := lottery_service.GetLuckyForPlayer(language)
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorGetLuckytodyFail, nil)
		return
	}
	resp := GetLotteryForUserResponse{LotteryContent: *lottery, LuckyContent: *luckyTody}
	appG.Response(http.StatusOK, e.SUCCESS, resp)
}
