package player

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/service/player/lottery_service"
)

// GetLottery
// @Summary 获取日签
// @Produce  json
// @Param uid query int true "用户uid"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /player/lottery [get]
// @Tags Player

func GetLottery(c *gin.Context) {
	appG := app.Gin{C: c}
	uid, _ := c.GetQuery("uid")
	fmt.Printf("uid: %s", uid)
	logging.Debug(uid)
	lottery, err := lottery_service.GetLottery()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
	}
	appG.Response(http.StatusOK, e.SUCCESS, lottery)
}
