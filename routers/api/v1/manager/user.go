package manager

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/service/manager/user_service"
)

type AddUserForm struct {
	Name   string `form:"name" binding:"required"`
	Passwd string `form:"passwd" binding:"required"`
}

// AddUser
// @Summary 添加用户
// @Param name formData string true "name"
// @Param passwd formData string true "passwd"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/user [post]
// @Tags manager
// @Security ApiKeyAuth
func AddUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		user AddUserForm
	)

	if err := c.ShouldBind(&user); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	userService := user_service.User{
		Name:   user.Name,
		Passwd: user.Passwd,
	}
	exists, err := userService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_USER, nil)
		return
	}

	if exists {
		appG.Response(http.StatusOK, e.ERROR_USER_HAS_EXIST, nil)
		return
	}
	err = userService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CREAT_USER, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

// GetUser
// @Summary 获取用户
// @Param id query int false "id"
// @Param name query string false "name"
// @Router /manager/user [get]
// @Security ApiKeyAuth
// @Tags manager
// @Success 200 {object} GetUserResponse
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
func GetUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		user GetUserForm
	)

	if err := c.ShouldBind(&user); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
}

type GetUserForm struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type GetUserResponse struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

// GetCurrentLoginUserInfo
// @Summary 通过登录态获取当前登录用户信息
// @Router /manager/user/info [get]
// @Security ApiKeyAuth
// @Tags manager
// @Success 200 {object} GetUserResponse
// @Failure 500 {object} app.Response
func GetCurrentLoginUserInfo(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	// TODO 需要还原
	username := c.GetString("username")

	userService := user_service.User{
		Name: username,
	}
	user, err := userService.GetUserByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	var resp GetUserResponse
	resp.Id = user.ID
	resp.Name = user.Name
	appG.Response(http.StatusOK, e.SUCCESS, resp)

}
