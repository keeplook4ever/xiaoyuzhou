package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/util"
	"xiaoyuzhou/service/lucky_service"
)

type AddLuckyForm struct {
	Data     []string `json:"data" binding:"required"`                         //字符串数组
	Type     string   `json:"type" binding:"required" enums:"spell,song,todo"` //咒语：spell, 歌曲：song, 适宜：todo
	Language string   `json:"language" binding:"required" enums:"jp,zh,en,tc"` //语言
}

// AddLucky
// @Summary 添加今日好运内容咒语
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
	luckyI := lucky_service.LuckyInputContent{
		Lists:    reqData.Data,
		Type:     reqData.Type,
		Language: reqData.Language,
	}
	if err := luckyI.Add(); err != nil {
		appG.Response(http.StatusOK, err.Error(), nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// UploadLucky
// @Summary 上传今日好运excel
// @Param file formData file true "excel文件"
// @Accept json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lucky/upload [post]
// @Security ApiKeyAuth
// @Tags Manager
func UploadLucky(c *gin.Context) {
	appG := app.Gin{C: c}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		appG.Response(http.StatusBadRequest, "参数错误", nil)
		return
	}
	if err := lucky_service.Import(file); err != nil {
		appG.Response(http.StatusOK, "导入excel失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// EditLucky
// @Summary 修改今日好运内容
// @Param _ body EditLuckyForm true "修改好运内容"
// @Produce json
// @Accept json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/lucky [put]
// @Security ApiKeyAuth
// @Tags Manager
func EditLucky(c *gin.Context) {
	appG := app.Gin{C: c}
	var edF EditLuckyForm
	if err := c.ShouldBindJSON(&edF); err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	if err := lucky_service.EditLucky(edF.Type, edF.Id, edF.Content, edF.Language); err != nil {
		appG.Response(http.StatusOK, "编辑失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//DeleteLucky
//@Summary 删除今日好运内容
//@Param _ body DeleteLuckyForm  true "删除lucky"
//@Produce json
//@Accept json
//@Success 200 {object} app.Response
//@Failure 500 {object} app.Response
//@Router /manager/lucky [delete]
//@Security ApiKeyAuth
//@Tags Manager
func DeleteLucky(c *gin.Context) {
	appG := app.Gin{C: c}
	var delF DeleteLuckyForm
	if err := c.ShouldBindJSON(&delF); err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	if err := lucky_service.Delete(delF.Type, delF.IdList); err != nil {
		appG.Response(http.StatusOK, "删除失败", nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// GetLucky
// @Summary 获取今日好运内容
// @Param type query string true "咒语\歌曲\适宜" Enums(spell,song,todo)
// @Param language query string true "语言" Enums(jp,zh,en,tc)
// @Produce json
// @Success 200 {object} GetLuckyResponse
// @Failure 500 {object} app.Response
// @Router /manager/lucky [get]
// @Security ApiKeyAuth
// @Tags Manager
func GetLucky(c *gin.Context) {
	appG := app.Gin{C: c}
	luckI := lucky_service.LuckyInputContent{
		Type:     c.Query("type"),
		Language: c.Query("language"),
		PageNum:  util.GetPage(c),
		PageSize: util.GetPageSize(c),
	}
	_type, lucks, count, err := luckI.Get()
	if err != nil {
		appG.Response(http.StatusOK, err.Error(), nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, GetLuckyResponse{
		Lists: lucks,
		Count: count,
		Type:  _type,
	})
}

type GetLuckyResponse struct {
	Lists interface{} `json:"lists"`
	Count int64       `json:"count"`
	Type  string      `json:"type"`
}

type DeleteLuckyForm struct {
	Type   string `json:"type" binding:"required" enums:"spell,song,todo"` //spell:咒语，song:歌曲, todo:适宜
	IdList []int  `json:"id_list" binding:"required"`
}

type EditLuckyForm struct {
	Id       int    `json:"id" binding:"required"`
	Type     string `json:"type" binding:"required" enums:"spell,song,todo"`
	Content  string `json:"content" binding:"required"`
	Language string `json:"language" binding:"required" enums:"jp,zh,en,tc"`
}
