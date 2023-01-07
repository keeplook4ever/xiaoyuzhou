package manager

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/aws"
	"xiaoyuzhou/pkg/e"
)

// GetS3Token
// @Summary 获取S3上传Token
// @Produce  json
// @Success 200 {object} aws.TmpTokenStruct
// @Failure 500 {object} app.Response
// @Tags Manager
// @Security ApiKeyAuth
// @Router /manager/s3/token [post]
func GetS3Token(c *gin.Context) {
	appG := app.Gin{C: c}
	var token *aws.TmpTokenStruct
	token, err := aws.GetToken()
	if err != nil {
		appG.Response(http.StatusInternalServerError, http.StatusInternalServerError, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, *token)

}
