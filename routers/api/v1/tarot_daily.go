package v1

// 每日免费塔罗

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/util"
	"xiaoyuzhou/service/tarot_service"
)

// GetDailyFreeTarot
// @Summary 获取每日免费塔罗牌
// @Param uid query string true "用户ID"
// @Param language query string true "语言"
// @Success 200 {object} models.TarotDto
// @Failure 500 {object} app.Response
// @Router /player/tarot-daily [get]
// @Tags Player
func GetDailyFreeTarot(c *gin.Context) {
	appG := app.Gin{C: c}
	tarot, err := tarot_service.GetDailyTarotFree(c.Query("uid"), c.Query("language"))
	if err != nil {
		appG.Response(http.StatusOK, "每日塔罗获取失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, tarot)
}

// AddDailyTarot
// @Summary 添加每日免费塔罗牌内容
// @Param _ body AddDailyTarotForm true "参数"
// @Produce json
// @Accept json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/tarot-daily [post]
// @Security ApiKeyAuth
// @Tags Manager
func AddDailyTarot(c *gin.Context) {
	var (
		data AddDailyTarotForm
		appG = app.Gin{C: c}
	)
	if err := c.ShouldBindJSON(&data); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	articleService := tarot_service.DailyTarotInput{
		ImgUrl:    data.ImgUrl,
		Language:  data.Language,
		CardName:  data.CardName,
		Analyze:   data.Analyze,
		LoveList:  data.LoveList,
		WorkList:  data.WorkList,
		CreatedBy: c.GetString("username"),
		UpdatedBy: c.GetString("username"),
	}
	if err := articleService.Add(); err != nil {
		appG.Response(http.StatusOK, err.Error(), nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type AddDailyTarotForm struct {
	ImgUrl    string   `json:"img_url" binding:"required"`                   // 图片链接
	Language  string   `json:"language" enums:"jp,zh,en" binding:"required"` // 语言
	CardName  string   `json:"card_name" binding:"required"`                 // 卡牌名字
	Analyze   string   `json:"analyze" binding:"required"`                   // 解读
	LoveList  []string `json:"love_list" binding:"required"`                 // 爱情列表
	WorkList  []string `json:"work_list" binding:"required"`                 // 工作列表
	CreatedBy string   `json:"created_by" binding:"required"`                // 创建者
	UpdatedBy string   `json:"updated_by" binding:"required"`                // 更新者
}

// EditDailyTarot
// @Summary 修改每日塔罗牌
// @Param _ body EditDailyTarotForm true "修改内容"
// @Param id path int true "塔罗牌ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Tags Manager
// @Security ApiKeyAuth
// @Router /manager/tarot-daily [put]
func EditDailyTarot(c *gin.Context) {
	var (
		data = EditDailyTarotForm{
			Id:        com.StrTo(c.Param("id")).MustInt(),
			UpdatedBy: c.GetString("username"),
		}
		appG = app.Gin{C: c}
	)
	if err := c.ShouldBindJSON(&data); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}
	tarotE := tarot_service.DailyTarotInput{
		ImgUrl:    data.ImgUrl,
		Language:  data.Language,
		CardName:  data.CardName,
		Analyze:   data.Analyze,
		LoveList:  data.LoveList,
		WorkList:  data.WorkList,
		UpdatedBy: c.GetString("username"),
	}

	exists, err := tarotE.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, "校验存在失败", nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, "不存在此塔罗牌", nil)
		return
	}
	err = tarotE.Edit()
	if err != nil {
		appG.Response(http.StatusOK, "编辑失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditDailyTarotForm struct {
	Id        int      `json:"id"`                        // 塔罗ID
	ImgUrl    string   `json:"img_url"`                   // 图片链接
	Language  string   `json:"language" enums:"jp,zh,en"` // 语言
	CardName  string   `json:"card_name"`                 // 卡牌名字
	Analyze   string   `json:"analyze"`                   // 解读
	LoveList  []string `json:"love_list"`                 // 爱情列表
	WorkList  []string `json:"work_list"`                 // 工作列表
	CreatedBy string   `json:"created_by"`                // 创建者
	UpdatedBy string   `json:"updated_by"`                // 更新者
}

// GetDailyTarot
// @Summary 获取每日塔罗牌:支持分页
// @Param id query int false "ID"
// @Param name query string false "名字"
// @Param language query string false "语言"
// @Success 200 {object} GetDailyTarotResponse
// @Failure 500 {object} app.Response
// @Router /manager/tarot-daily [get]
// @Security ApiKeyAuth
// @Tags Manager
func GetDailyTarot(c *gin.Context) {
	appG := app.Gin{C: c}
	tarotD := tarot_service.DailyTarotInput{
		Id:       com.StrTo(c.Query("id")).MustInt(),
		CardName: c.Query("name"),
		Language: c.Query("language"),
		PageNum:  util.GetPage(c),
		PageSize: util.GetPageSize(c),
	}
	tarots, count, err := tarotD.Get()
	if err != nil {
		appG.Response(http.StatusOK, "获取失败", nil)
		return
	}
	var resp GetDailyTarotResponse
	resp.Count = count
	resp.Lists = tarots
	appG.Response(http.StatusOK, e.SUCCESS, resp)
}

type GetDailyTarotResponse struct {
	Count int64                  `json:"count"`
	Lists []models.DailyTarotDto `json:"lists"`
}
