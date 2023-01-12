package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/util"
	"xiaoyuzhou/service/user_service"
)

type AddUserForm struct {
	Name   string `form:"name" binding:"required"`
	Passwd string `form:"passwd" binding:"required"`
	Role   string `form:"role" binding:"required"`
}

// AddUser
// @Summary 添加用户
// @Param name formData string true "name"
// @Param passwd formData string true "passwd"
// @Param role formData string true "role"
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
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	userService := user_service.UserInput{
		Name:      user.Name,
		Passwd:    user.Passwd,
		CreatedBy: c.GetString("username"),
		UpdatedBy: c.GetString("username"),
		Role:      user.Role,
	}
	exists, err := userService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorCheckExistUser, nil)
		return
	}

	if exists {
		appG.Response(http.StatusOK, e.ErrorUserHasExist, nil)
		return
	}
	err = userService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorCreatUser, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

// GetUser
// @Summary 获取用户
// @Param id query int false "id"
// @Param name query string false "name"
// @Param role query string false "role"
// @Router /manager/user [get]
// @Security ApiKeyAuth
// @Tags Manager
// @Success 200 {object} GetUserResponse
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
func GetUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	name := c.Query("name")
	role := c.Query("role")
	id := com.StrTo(c.Query("id")).MustInt()
	userService := user_service.UserInput{
		Name: name,
		ID:   uint(id),
		Role: role,
	}
	users, err := userService.GetUser()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorGetUserFail, nil)
		return
	}
	count := len(users)
	// 定义返回数据
	appG.Response(http.StatusOK, e.SUCCESS, GetUserResponse{Lists: users, Count: count})
}

type GetUserResponse struct {
	Lists []models.UserDto `json:"lists"`
	Count int              `json:"count"`
}

// GetCurrentLoginUserInfo
// @Summary 通过登录态获取当前登录用户信息
// @Router /manager/user/info [get]
// @Security ApiKeyAuth
// @Tags Manager
// @Success 200 {object} models.UserDto
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
		appG.Response(http.StatusOK, e.ErrorGetUserInfoFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, user[0])
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
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	userService := user_service.UserInput{Name: auth_.Username, Passwd: auth_.Password}
	isExist, err := userService.Check()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorAuthCheckTokenFail, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ErrorAuth, nil)
		return
	}

	token, err := util.GenerateToken(auth_.Username, auth_.Password)
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorAuthToken, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}

type TokenResponse struct {
	Token string `json:"token"`
}
