package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/util"
	"xiaoyuzhou/service/order_service"
	"xiaoyuzhou/service/tarot_service"
)

type AddTarotForm struct {
	ImgUrl        string   `json:"img_url" binding:"required"`                      // 图片链接
	Language      string   `json:"language" enums:"jp,zh,en,tc" binding:"required"` // 语言
	Pos           string   `json:"pos" enums:"up,down" binding:"required"`          // 塔罗正逆位: up正，down逆
	CardName      string   `json:"card_name" binding:"required"`                    // 卡牌名字
	KeyWord       string   `json:"keyword" binding:"required"`                      // 卡牌解读关键词
	Constellation string   `json:"constellation" binding:"required"`                // 对应星座
	People        string   `json:"people" binding:"required"`                       // 对应人物
	Element       string   `json:"element" binding:"required"`                      // 对应元素
	Enhance       string   `json:"enhance" binding:"required"`                      // 加强牌
	AnalyzeOne    string   `json:"analyze_one" binding:"required"`                  // 解析1
	AnalyzeTwo    string   `json:"analyze_two" binding:"required"`                  // 解析2
	PosMeaning    string   `json:"pos_meaning" binding:"required"`                  // 正逆位含义
	Love          string   `json:"love" binding:"required"`                         // 爱情婚姻
	Work          string   `json:"work" binding:"required"`                         // 事业学业
	Money         string   `json:"money" binding:"required"`                        // 人际财富
	Health        string   `json:"health" binding:"required"`                       // 健康生活
	Other         string   `json:"other" binding:"required"`                        // 其他
	LuckyNumber   string   `json:"lucky_number" binding:"required"`                 // 幸运数字
	Saying        string   `json:"saying" binding:"required"`                       // 名言
	AnswerList    []string `json:"answer_list" binding:"required"`                  // 答案列表
}

type EditTarotForm struct {
	Id            int      `json:"id"`                           // 塔罗牌ID
	ImgUrl        string   `json:"img_url"`                      // 图片链接
	Language      string   `json:"language" enums:"jp,zh,en,tc"` // 语言
	Pos           string   `json:"pos" enums:"up,down"`          // 塔罗正逆位 up正，down逆
	CardName      string   `json:"card_name"`                    // 卡牌名字
	KeyWord       string   `json:"keyword"`                      // 卡牌解读关键词
	Constellation string   `json:"constellation"`                // 对应星座
	People        string   `json:"people"`                       // 对应人物
	Element       string   `json:"element"`                      // 对应元素
	Enhance       string   `json:"enhance"`                      // 加强牌
	AnalyzeOne    string   `json:"analyze_one"`                  // 解析1
	AnalyzeTwo    string   `json:"analyze_two"`                  // 解析2
	PosMeaning    string   `json:"pos_meaning"`                  // 正逆位含义
	Love          string   `json:"love"`                         // 爱情婚姻
	Work          string   `json:"work"`                         // 事业学业
	Money         string   `json:"money"`                        // 人际财富
	Health        string   `json:"health"`                       // 健康生活
	Other         string   `json:"other"`                        // 其他
	LuckyNumber   string   `json:"lucky_number"`                 // 幸运数字
	Saying        string   `json:"saying"`                       // 名言
	AnswerList    []string `json:"answer_list"`                  // 答案列表
	UpdatedBy     string   `json:"updated_by"`                   // 修改人

}

// AddTarot
// @Summary 添加塔罗牌内容
// @Param _ body AddTarotForm true "参数"
// @Produce json
// @Accept json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/tarot [post]
// @Security ApiKeyAuth
// @Tags Manager
func AddTarot(c *gin.Context) {
	var (
		data AddTarotForm
		appG = app.Gin{C: c}
	)
	if err := c.ShouldBindJSON(&data); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	articleService := tarot_service.TarotInput{
		ImgUrl:        data.ImgUrl,
		Language:      data.Language,
		Pos:           data.Pos,
		CardName:      data.CardName,
		KeyWord:       data.KeyWord,
		Constellation: data.Constellation,
		People:        data.People,
		Element:       data.Element,
		Enhance:       data.Enhance,
		AnalyzeOne:    data.AnalyzeOne,
		AnalyzeTwo:    data.AnalyzeTwo,
		PosMeaning:    data.PosMeaning,
		Love:          data.Love,
		Work:          data.Work,
		Money:         data.Money,
		Health:        data.Health,
		Other:         data.Other,
		LuckyNumber:   data.LuckyNumber,
		Saying:        data.Saying,
		AnswerList:    data.AnswerList,
		CreatedBy:     c.GetString("username"),
		UpdatedBy:     c.GetString("username"),
	}
	if err := articleService.Add(); err != nil {
		appG.Response(http.StatusOK, err.Error(), nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// UploadTarot
// @Summary 上传塔罗牌excel, 批量创建(除了图片)
// @Param file formData file true "excel文件"
// @Accept json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/tarot/upload [post]
// @Security ApiKeyAuth
// @Tags Manager
func UploadTarot(c *gin.Context) {
	appG := app.Gin{C: c}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		appG.Response(http.StatusBadRequest, "参数错误", nil)
		return
	}
	if err := tarot_service.Import(file); err != nil {
		appG.Response(http.StatusOK, "导入塔罗excel失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// EditTarot
// @Summary 修改塔罗牌内容
// @Param _ body EditTarotForm true "参数"
// @Param id path int true "塔罗牌ID"
// @Produce json
// @Accept json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/tarot/{id} [put]
// @Security ApiKeyAuth
// @Tags Manager
func EditTarot(c *gin.Context) {
	var (
		data = EditTarotForm{
			Id:        com.StrTo(c.Param("id")).MustInt(),
			UpdatedBy: c.GetString("username"),
		}
		appG = app.Gin{C: c}
	)
	if err := c.ShouldBindJSON(&data); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	tarotS := tarot_service.TarotInput{
		ID:            data.Id,
		ImgUrl:        data.ImgUrl,
		Language:      data.Language,
		Pos:           data.Pos,
		CardName:      data.CardName,
		KeyWord:       data.KeyWord,
		Constellation: data.Constellation,
		People:        data.People,
		Element:       data.Element,
		Enhance:       data.Enhance,
		AnalyzeOne:    data.AnalyzeOne,
		AnalyzeTwo:    data.AnalyzeTwo,
		PosMeaning:    data.PosMeaning,
		Love:          data.Love,
		Work:          data.Work,
		Money:         data.Money,
		Health:        data.Health,
		Other:         data.Other,
		LuckyNumber:   data.LuckyNumber,
		Saying:        data.Saying,
		AnswerList:    data.AnswerList,
		UpdatedBy:     c.GetString("username"),
	}

	exists, err := tarotS.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, "校验存在失败", nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, "不存在此塔罗牌", nil)
		return
	}
	if err := tarotS.Edit(); err != nil {
		appG.Response(http.StatusOK, "编辑失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type GetTarotResponse struct {
	Count int64             `json:"count"` // 总数
	Lists []models.TarotDto `json:"lists"` // 塔罗牌列表
}

// GetTarot
// @Summary 获取塔罗牌:支持分页
// @Param id query int false "ID"
// @Param name query string false "名字"
// @Param language query string false "语言"
// @Param pos query string false "正逆位"
// @Param keyword query string false "关键字"
// @Param constellation query string false "星座"
// @Param people query string false "对应人物"
// @Param element query string false "对应元素"
// @Param enhance query string false "加强牌"
// @Success 200 {object} GetTarotResponse
// @Failure 500 {object} app.Response
// @Router /manager/tarot [get]
// @Security ApiKeyAuth
// @Tags Manager
func GetTarot(c *gin.Context) {
	appG := app.Gin{C: c}
	tarotS := tarot_service.TarotInput{
		ID:            com.StrTo(c.Query("id")).MustInt(),
		CardName:      c.Query("name"),
		Language:      c.Query("language"),
		Pos:           c.Query("pos"),
		KeyWord:       c.Query("keyword"),
		Constellation: c.Query("constellation"),
		People:        c.Query("pepole"),
		Element:       c.Query("element"),
		Enhance:       c.Query("enhance"),
		PageNum:       util.GetPage(c),
		PageSize:      util.GetPageSize(c),
	}
	tarots, count, err := tarotS.Get()
	if err != nil {
		appG.Response(http.StatusOK, "获取失败", nil)
		return
	}
	var resp GetTarotResponse
	resp.Count = count
	resp.Lists = tarots
	appG.Response(http.StatusOK, e.SUCCESS, resp)
}

// GetTarotOne
// @Summary 用户抽取塔罗牌: 1张
// @Param _ body GetTarotOneForm true "抽取一张塔罗牌付费"
// @Success 200 {object} GetTarotOneRes
// @Failure 500 {object} app.Response
// @Router /player/tarot/one [post]
// @Tags Player
func GetTarotOne(c *gin.Context) {
	appG := app.Gin{C: c}
	var formD GetTarotOneForm
	if err := c.ShouldBindJSON(&formD); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	tarot, orderId, err := tarot_service.GetRandomOneTarot(formD.Uid, formD.Question, formD.Language)
	if err != nil {
		appG.Response(http.StatusOK, "获取失败", nil)
		return
	}
	var resp GetTarotOneRes
	resp.TarotID = tarot.TarotId
	resp.Name = tarot.CardName
	resp.ImgUrl = tarot.ImgUrl
	resp.OrderId = orderId
	data, err := models.GetPrice(formD.Language)
	if err != nil {
		appG.Response(http.StatusOK, "获取价格失败", nil)
		return
	}
	resp.Price = *data
	appG.Response(http.StatusOK, e.SUCCESS, resp)
}

// GetTarotAnswer
// @Summary 获取用户抽取塔罗牌的解答
// @Param order_id query string true "订单ID"
// @Success 200 {object} GetTarotOneAnswerResp
// @Failure 500 {object} app.Response
// @Router /player/tarot/answer [get]
// @Tags Player
func GetTarotAnswer(c *gin.Context) {
	appG := app.Gin{C: c}
	payed, err := order_service.CheckOrderIfPayed(c.Query("order_id"))
	if err != nil {
		appG.Response(http.StatusOK, "订单校验失败", nil)
		return
	}
	if !payed {
		appG.Response(http.StatusOK, err.Error(), nil)
		return
	}
	// 根据订单号查对应抽取塔罗牌的答案, 这里的订单号是原始订单号
	tarot, question, ts, err := tarot_service.GetOneTarotByOrderAndUser(c.Query("order_id"))
	if err != nil {
		appG.Response(http.StatusOK, "获取塔罗失败", nil)
		return
	}

	var resp GetTarotOneAnswerResp
	resp.Type = "one"
	resp.Tarot = *tarot
	resp.Question = question
	resp.Ts = ts
	appG.Response(http.StatusOK, e.SUCCESS, resp)
}

type GetTarotOneRes struct {
	TarotID uint         `json:"tarot_id"` // 塔罗牌id
	Name    string       `json:"name"`     // 塔罗牌名字
	ImgUrl  string       `json:"img_url"`  // 塔罗牌图片链接
	Price   models.Price `json:"price"`    // 价格
	OrderId string       `json:"order_id"` // 用户抽塔罗牌的订单记录ID
}

type GetTarotOneForm struct {
	Uid      string `json:"uid" binding:"required"`      // 用户uid
	Question string `json:"question" binding:"required"` // 用户问题
	Language string `json:"language" binding:"required" enums:"jp,zh,en,tc"`
}

type GetTarotOneAnswerResp struct {
	Tarot    models.TarotDto `json:"tarot"`                  // 塔罗详情
	Question string          `json:"question"`               // 用户提的问题
	Ts       int64           `json:"ts"`                     // 用户抽取塔罗牌的时间戳
	Type     string          `json:"type" enums:"one,three"` // 单张或多张
}
