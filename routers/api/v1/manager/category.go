package manager

import (
	"net/http"
	"xiaoyuzhou/service/manager/category_service"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/export"
	"xiaoyuzhou/pkg/setting"
	"xiaoyuzhou/pkg/util"
)

// GetTags
// @Summary 获取文章类型
// @Accept json
// @Produce  json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Tags Manager
// @Security ApiKeyAuth
// @Router /manager/tags [get]
func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("name")

	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := category_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	tags, err := tagService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_CATEGORYS_FAIL, nil)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_CATEGORY_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})
}

type AddTagForm struct {
	Name      string `json:"name" binding:"required"`
	CreatedBy string `json:"created_by" binding:"required"`
	State     int    `json:"state" binding:"required"`
}

// AddTag
// @Summary 添加文章类型
// @Produce  json
// @Param _ body AddTagForm true "文章类型"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/tags [post]
// @Tags Manager
// @Security ApiKeyAuth
func AddTag(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddTagForm
	)
	if err := c.ShouldBindJSON(&form); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	tagService := category_service.Tag{
		Name:      form.Name,
		CreatedBy: form.CreatedBy,
		State:     form.State,
	}
	exists, err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_CATEGORY_FAIL, nil)
		return
	}
	if exists {
		appG.Response(http.StatusOK, e.ERROR_EXIST_CATEGORY, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_CATEGORY_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditTagForm struct {
	ID         int    `form:"id" valid:"required;Min(1)"`
	Name       string `form:"name" valid:"required;MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0,1)"`
}

// EditTag
// @Summary 修改文章类型
// @Produce  json
// @Param id path int true "ID"
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param modified_by body string true "ModifiedBy"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/tags/{id} [put]
// @Tags Manager
// @Security ApiKeyAuth
func EditTag(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form = EditTagForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	tagService := category_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}

	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_CATEGORY_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_CATEGORY, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_CATEGORY_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// DeleteTag
// @Summary 删除文章类型
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/tags/{id} [delete]
// @Tags Manager
// @Security ApiKeyAuth
func DeleteTag(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	}

	tagService := category_service.Tag{ID: id}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_CATEGORY_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_CATEGORY, nil)
		return
	}

	if err := tagService.Delete(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_CATEGORY_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// ExportTag
// @Summary 导出文章类型
// @Produce  json
// @Param name body string false "Name"
// @Param state body int false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/tags/export [post]
// @Tags Manager
// @Security ApiKeyAuth
func ExportTag(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.PostForm("name")
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := category_service.Tag{
		Name:  name,
		State: state,
	}

	filename, err := tagService.Export()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXPORT_CATEGORY_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelPath() + filename,
	})
}
