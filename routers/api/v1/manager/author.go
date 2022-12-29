package manager

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/service/manager/author_service"
)

type AddAuthorForm struct {
	Name      string `json:"name" binding:"required"`
	Age       int    `json:"age" binding:"required"`
	Gender    int    `json:"gender" binding:"required"`
	CreatedBy string `json:"createdBy" binding:"required"`
	Desc      string `json:"desc" binding:"required"`
}

type EditAuthorForm struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Gender     int    `json:"gender"`
	ModifiedBy string `json:"modifiedBy" binding:"required"`
	Desc       string `json:"desc"`
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
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	authorService := author_service.Author{
		Name:      author.Name,
		Age:       author.Age,
		Gender:    author.Gender,
		Desc:      author.Desc,
		CreatedBy: author.CreatedBy,
	}
	exists, err := authorService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_AUTHOR_FAIL, nil)
		return
	}
	if exists {
		appG.Response(http.StatusOK, e.ERROR_EXIST_AUTHOR, nil)
		return
	}

	err = authorService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_AUTHOR_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// EditAuthor
// @Summary 编辑作者
// @Produce json
// @Param id path int true "ID"
// @Param _ body EditAuthorForm false "编辑作者"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/author/{id} [put]
// @Tags Manager
// @Security ApiKeyAuth
func EditAuthor(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		author EditAuthorForm
	)

	if err := c.ShouldBindJSON(&author); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
}
