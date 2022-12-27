package player

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/logging"
)

// GetArticleForPlayer
// @Summary "为玩家展示文章"
// @Param uid query int true "用户id"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /player/article [get]
// @Tags Player
func GetArticleForPlayer(c *gin.Context) {
	uid := c.Query("uid")
	appG := app.Gin{C: c}
	logging.Debugf("uid: ", uid)
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
