package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/pkg/util"
	"xiaoyuzhou/service/trueword_service"
)

type AddTrueWordForm struct {
	Lang     string   `json:"lang" binding:"required" enums:"jp,zh,en,ts"` // 语言
	WordList []string `json:"word_list" binding:"required"`                // 真言数组
}

type EditTrueWordForm struct {
	Lang string `json:"lang" enums:"jp,zh,en,ts"` // 语言
	Word string `json:"word" binding:"required"`  // 真言
}

type GetTrueWordResponse struct {
	Count int64             `json:"count"`
	Lists []models.TrueWord `json:"lists"`
}

// AddTrueWord
// @Summary 添加真言
// @Param _ body AddTrueWordForm true "参数"
// @Produce json
// @Accept json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/true-word [post]
// @Security ApiKeyAuth
// @Tags Manager
func AddTrueWord(c *gin.Context) {
	appG := app.Gin{C: c}
	var _data AddTrueWordForm
	if err := c.ShouldBindJSON(&_data); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}
	trueWord := trueword_service.TrueWordInput{
		Lang:      _data.Lang,
		WordList:  _data.WordList,
		UpdatedBy: c.GetString("username"),
		CreatedBy: c.GetString("username"),
	}
	if err := trueWord.Add(); err != nil {
		appG.Response(http.StatusOK, "添加失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// EditTrueWord
// @Summary 编辑真言:一次只能编辑一条
// @Param _ body EditTrueWordForm true "编辑body"
// @Param id path int true "id"
// @Produce json
// @Accept json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/true-word/{id} [put]
// @Security ApiKeyAuth
// @Tags Manager
func EditTrueWord(c *gin.Context) {
	appG := app.Gin{C: c}
	var _data EditTrueWordForm
	if err := c.ShouldBindJSON(&_data); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}
	//一次只能编辑一条
	wordList := make([]string, 0)
	wordList = append(wordList, _data.Word)
	trueWord := trueword_service.TrueWordInput{
		Id:        com.StrTo(c.Param("id")).MustInt(),
		WordList:  wordList,
		Lang:      _data.Lang,
		UpdatedBy: c.GetString("username"),
	}
	if err := trueWord.Edit(); err != nil {
		appG.Response(http.StatusOK, err.Error(), nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// DeleteTrueWord
// @Summary 删除真言
// @Param id path int true "id"
// @Produce json
// @Accept json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/true-word/{id} [delete]
// @Security ApiKeyAuth
// @Tags Manager
func DeleteTrueWord(c *gin.Context) {
	appG := app.Gin{C: c}
	trueWord := trueword_service.TrueWordInput{
		Id: com.StrTo(c.Param("id")).MustInt(),
	}
	if err := trueWord.Delete(); err != nil {
		appG.Response(http.StatusOK, err.Error(), nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// GetTrueWord
// @Summary 获取真言
// @Param id query int false "编辑body"
// @Param language query string false "语言类型" Enums(en,jp,zh,ts)
// @Produce json
// @Accept json
// @Success 200 {object} GetTrueWordResponse
// @Failure 500 {object} app.Response
// @Router /manager/true-word [get]
// @Security ApiKeyAuth
// @Tags Manager
func GetTrueWord(c *gin.Context) {
	appG := app.Gin{C: c}
	trueWord := trueword_service.TrueWordInput{
		Id:       com.StrTo(c.Query("id")).MustInt(),
		Lang:     c.Query("language"),
		PageNum:  util.GetPage(c),
		PageSize: util.GetPageSize(c),
	}
	twL, count, err := trueWord.Get()
	if err != nil {
		appG.Response(http.StatusOK, err.Error(), nil)
		return
	}
	var res GetTrueWordResponse
	res.Lists = twL
	res.Count = count
	appG.Response(http.StatusOK, e.SUCCESS, res)
}

// GetTrueWordForPlayer
// @Summary 为用户输出真言
// @Param language query string false "语言类型" Enums(en,jp,zh,ts)
// @Produce json
// @Accept json
// @Success 200 {object} models.TrueWord
// @Failure 500 {object} app.Response
// @Router /player/true-word [get]
// @Tags Player
func GetTrueWordForPlayer(c *gin.Context) {
	appG := app.Gin{C: c}
	trueWord := trueword_service.TrueWordInput{
		Lang:     c.Query("language"),
		PageNum:  util.GetPage(c),
		PageSize: util.GetPageSize(c),
	}
	twL, count, err := trueWord.Get()
	if err != nil {
		appG.Response(http.StatusOK, err.Error(), nil)
		return
	}
	// 从列表中获取随机的一个给前端
	if count == 0 {
		logging.Error("没有真言了！需要后台创建！")
		appG.Response(http.StatusOK, "没有真言,需创建", nil)
		return
	}
	ix := util.RandFromRange(0, int(count))
	tw := twL[ix]
	appG.Response(http.StatusOK, e.SUCCESS, tw)
}
