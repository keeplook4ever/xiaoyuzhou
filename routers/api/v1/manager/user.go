package manager

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaoyuzhou/models/manager"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/util"
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
// @Tags Manager
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

	userService := user_service.UserInput{
		Name:      user.Name,
		Passwd:    user.Passwd,
		CreatedBy: c.GetString("username"),
		UpdatedBy: c.GetString("username"),
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
// @Tags Manager
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

	userService := user_service.UserInput{}
	users, err := userService.GetUser()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	count := len(users)
	// 定义返回数据
	appG.Response(http.StatusOK, e.SUCCESS, GetUserResponse{Lists: users, Count: count})
}

type GetUserForm struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type GetUserResponse struct {
	Lists []manager.UserDto `json:"lists"`
	Count int               `json:"count"`
}

// GetCurrentLoginUserInfo
// @Summary 通过登录态获取当前登录用户信息
// @Router /manager/user/info [get]
// @Security ApiKeyAuth
// @Tags Manager
// @Success 200 {object} GetUserResponse
// @Failure 500 {object} app.Response
func GetCurrentLoginUserInfo(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	username := c.GetString("username")

	userService := user_service.UserInput{
		Name: username,
	}
	user, err := userService.GetUser()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, GetUserResponse{
		Lists: user,
		Count: 1,
	})
}

type auth struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// GetAuth
// @Summary 获取Token
// @Accept json
// @Produce  json
// @Param _ body auth true "用户名和密码"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/user/auth [post]
// @Tags Manager
func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	var auth_ auth
	if err := c.ShouldBindJSON(&auth_); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	userService := user_service.UserInput{Name: auth_.Username, Passwd: auth_.Password}
	isExist, err := userService.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(auth_.Username, auth_.Password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}

type TokenResponse struct {
	Token string `json:"token"`
}
