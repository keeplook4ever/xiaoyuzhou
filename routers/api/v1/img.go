package v1

import (
	"net/http"
	"xiaoyuzhou/pkg/qiniu"
	"xiaoyuzhou/pkg/util"

	"github.com/gin-gonic/gin"

	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/pkg/upload"
)

// UploadImage
// @Summary 上传图片
// @Produce  json
// @Param image formData file true "Image File"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/articles/img [post]
// @Tags Manager
// @Security ApiKeyAuth
func UploadImage(c *gin.Context) {
	appG := app.Gin{C: c}
	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if image == nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}
	imageName := image.Filename
	fileMd5 := util.EncodeMD5(imageName)
	if !upload.CheckImageExt(imageName) {
		appG.Response(http.StatusBadRequest, e.ErrorUploadCheckImageFormat, nil)
		return
	}

	size := image.Size

	// 上传七牛云
	imgUrl, err := qiniu.UploadImg(file, size, fileMd5)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorPutFileToQiniu, nil)
		return
	}
	resp := make(map[string]interface{})

	resp["img_url"] = imgUrl
	appG.Response(http.StatusOK, e.SUCCESS, resp)
	return

}
