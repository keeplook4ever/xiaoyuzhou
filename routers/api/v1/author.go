package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/setting"
	"xiaoyuzhou/pkg/util"
	"xiaoyuzhou/service/author_service"
)

type AddAuthorForm struct {
	Name   string `json:"name" binding:"required"`
	Age    int    `json:"age" binding:"required"`
	Gender int    `json:"gender" binding:"required" enums:"1,2" default:"2"` //1表示男，2表示女
	Desc   string `json:"desc" binding:"required"`
}

// AddAuthor
// @Summary 添加作者
// @Produce  json
// @Param _ body AddAuthorForm true "作者详情"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/author [post]
// @Tags Manager
// @Security ApiKeyAuth
func AddAuthor(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		author AddAuthorForm
	)

	if err := c.ShouldBindJSON(&author); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	authorService := author_service.AuthorInput{
		Name:      author.Name,
		Age:       author.Age,
		Gender:    author.Gender,
		Desc:      author.Desc,
		CreatedBy: c.GetString("username"), // 创建者从登录用户token获取
		UpdatedBy: c.GetString("username"), //默认更新者是创建者
	}
	exists, err := authorService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorExistAuthorFail, nil)
		return
	}
	if exists {
		appG.Response(http.StatusOK, e.ErrorExistAuthor, nil)
		return
	}

	err = authorService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorAddAuthorFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditAuthorForm struct {
	Name   string `form:"name"`
	Age    int    `form:"age"`
	Gender int    `form:"gender" enums:"1,2" default:"2"` //1表示男，2表示女
	Desc   string `form:"desc" binding:"required"`
	Id     int    `form:"id" binding:"required"`
}

// EditAuthor
// @Summary 编辑作者
// @Produce json
// @Param id path int true "ID"
// @Param name formData string false "Name"
// @Param age formData int false "Age"
// @Param gender formData int false "Gender" Enums(1,2) default(2)
// @Param desc formData string true "Desc"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/author/{id} [put]
// @Tags Manager
// @Security ApiKeyAuth
func EditAuthor(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		author = EditAuthorForm{}
	)

	if err := c.ShouldBind(&author); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	authorService := author_service.AuthorInput{
		Name:      author.Name,
		Age:       author.Age,
		Gender:    author.Gender,
		Desc:      author.Desc,
		ID:        author.Id,
		UpdatedBy: c.GetString("username"), //修改者从登录用户态获取
	}

	exists, err := authorService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorExistAuthorFail, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ErrorNotExistAuthor, nil)
		return
	}

	err = authorService.Edit()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorEditAuthorFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type GetAuthorsResponse struct {
	Lists []models.AuthorDto `json:"lists"`
	Count int64              `json:"count"`
}

// GetAuthors
// @Summary 获取作者
// @Produce json
// @Param name query string false "Name"
// @Param id query string false "ID"
// @Success 200 {object} GetAuthorsResponse
// @Failure 500 {object} app.Response
// @Router /manager/author [get]
// @Tags Manager
// @Security ApiKeyAuth
func GetAuthors(c *gin.Context) {
	var appG = app.Gin{C: c}
	name := c.Query("name")
	id := com.StrTo(c.Query("id")).MustInt()
	authorService := author_service.AuthorInput{
		Name:     name,
		ID:       id,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	authors, err := authorService.GetAll()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorGetAuthorFail, nil)
		return
	}

	count, err := authorService.Count()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorCountAuthorFail, nil)
		return
	}

	var res GetAuthorsResponse
	res.Lists = authors
	res.Count = count

	appG.Response(http.StatusOK, e.SUCCESS, res)
}
