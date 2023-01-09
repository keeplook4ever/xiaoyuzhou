package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/service/lottery_service"
)

type GetLotteryResponse struct {
	LotteryContent models.LotteryDto
	LuckyContent   models.LuckyTodayDto
}

// GetLottery
// @Summary 获取日签
// @Produce  json
// @Param uid query int true "用户uid"
// @Success 200 {object} GetLotteryResponse
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /player/lottery [get]
// @Tags Player
func GetLottery(c *gin.Context) {
	appG := app.Gin{C: c}
	uid := c.Query("uid")
	logging.Debugf("uid is: %s", uid)
	// 存储用户uid的记录
	lottery, err := lottery_service.GetLottery()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorGetLotteryFail, nil)
		return
	}
	luckyTody, err := lottery_service.GetLucky()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorGetLuckytodyFail, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, GetLotteryResponse{LotteryContent: lottery, LuckyContent: luckyTody})
}
