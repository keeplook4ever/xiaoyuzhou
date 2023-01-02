package manager

import (
	"net/http"
	"xiaoyuzhou/models/manager"
	"xiaoyuzhou/service/manager/category_service"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/setting"
	"xiaoyuzhou/pkg/util"
)

type GetCategoryResponse struct {
	Lists []manager.Category `json:"lists"`
	Count int                `json:"count"`
}

// GetCategory
// @Summary 获取文章类型
// @Accept json
// @Produce  json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Param id query int false "ID"
// @Success 200 {object} GetCategoryResponse
// @Failure 500 {object} app.Response
// @Tags Manager
// @Security ApiKeyAuth
// @Router /manager/category [get]
func GetCategory(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("name")
	id := com.StrTo(c.Query("id")).MustInt()
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	categoryService := category_service.Category{
		ID:       id,
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	categories, err := categoryService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_CATEGORYS_FAIL, nil)
		return
	}

	count, err := categoryService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_CATEGORY_FAIL, nil)
		return
	}

	var res GetCategoryResponse
	res.Lists = categories
	res.Count = count

	appG.Response(http.StatusOK, e.SUCCESS, res)
}

type AddCategoryForm struct {
	Name  string `json:"name" binding:"required"`
	State int    `json:"state" binding:"required" default:"1"`
}

// AddCategory
// @Summary 添加文章类型
// @Produce  json
// @Param _ body AddCategoryForm true "文章类型"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/category [post]
// @Tags Manager
// @Security ApiKeyAuth
func AddCategory(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddCategoryForm
	)
	if err := c.ShouldBindJSON(&form); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	categoryService := category_service.Category{
		Name:      form.Name,
		CreatedBy: "", // 后端获取登录态用户
		State:     form.State,
	}
	exists, err := categoryService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_CATEGORY_FAIL, nil)
		return
	}
	if exists {
		appG.Response(http.StatusOK, e.ERROR_EXIST_CATEGORY, nil)
		return
	}

	err = categoryService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_CATEGORY_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditCategoryForm struct {
	ID    int    `form:"id" binding:"required"`
	Name  string `form:"name"`
	State int    `form:"state" enums:"0,1"` // 0表示禁用，1表示启用
}

// EditCategory
// @Summary 修改文章类型
// @Produce  json
// @Param id path int true "ID"
// @Param name formData string false "Name"
// @Param state formData int false "State" default(1)
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/category/{id} [put]
// @Tags Manager
// @Security ApiKeyAuth
func EditCategory(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form = EditCategoryForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)
	if err := c.ShouldBind(&form); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	categoryService := category_service.Category{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: "", // 修改者从登录用户态获取
		State:      form.State,
	}

	exists, err := categoryService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_CATEGORY_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_CATEGORY, nil)
		return
	}

	err = categoryService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_CATEGORY_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// DeleteCategory
// @Summary 删除文章类型
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/category/{id} [delete]
// @Tags Manager
// @Security ApiKeyAuth
func DeleteCategory(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	}

	categoryService := category_service.Category{ID: id}
	exists, err := categoryService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_CATEGORY_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_CATEGORY, nil)
		return
	}

	if err = categoryService.Delete(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_CATEGORY_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
