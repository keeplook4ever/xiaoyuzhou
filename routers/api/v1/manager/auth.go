package manager

import (
	"net/http"
	"xiaoyuzhou/service/manager/auth_service"

	"github.com/gin-gonic/gin"

	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/util"
)

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
// @Router /auth [post]
// @Tags Manager
func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	var auth_ auth
	if err := c.ShouldBindJSON(&auth_); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	authService := auth_service.Auth{Username: auth_.Username, Password: auth_.Password}
	isExist, err := authService.Check()
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
