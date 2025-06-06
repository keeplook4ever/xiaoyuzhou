package v1

import (
	"net/http"
	"xiaoyuzhou/models"
	"xiaoyuzhou/service/category_service"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/util"
)

type GetCategoryResponse struct {
	Lists []models.CategoryDto `json:"lists"`
	Count int64                `json:"count"`
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

	categoryService := category_service.CategoryInput{
		ID:       id,
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: util.GetPageSize(c),
	}
	categories, count, err := categoryService.GetAll()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorGetCategoriesFail, nil)
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
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	categoryService := category_service.CategoryInput{
		Name:      form.Name,
		CreatedBy: c.GetString("username"), // 后端获取登录态用户
		State:     form.State,
		UpdatedBy: c.GetString("username"), // 默认更新者是创建者
	}
	exists, err := categoryService.ExistByName()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorExistCategoryFail, nil)
		return
	}
	if exists {
		appG.Response(http.StatusOK, e.ErrorExistCategory, nil)
		return
	}

	err = categoryService.Add()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorAddCategoryFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditCategoryForm struct {
	ID    int    `json:"id" binding:"required"`
	Name  string `json:"name"`
	State int    `json:"state" default:"1"`
}

// EditCategory
// @Summary 修改文章类型
// @Produce  json
// @Param id path int true "ID"
// @Param _ body EditCategoryForm true "修改类型参数"
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
	if err := c.ShouldBindJSON(&form); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	categoryService := category_service.CategoryInput{
		ID:        form.ID,
		Name:      form.Name,
		UpdatedBy: c.GetString("username"), // 修改者从登录用户态获取
		State:     form.State,
	}

	exists, err := categoryService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorExistCategoryFail, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ErrorNotExistCategory, nil)
		return
	}

	err = categoryService.Edit()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorEditCategoryFail, nil)
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
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
	}

	categoryService := category_service.CategoryInput{ID: id}
	exists, err := categoryService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorExistCategoryFail, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ErrorNotExistCategory, nil)
		return
	}

	if err = categoryService.Delete(); err != nil {
		appG.Response(http.StatusOK, e.ErrorDeleteCategoryFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
