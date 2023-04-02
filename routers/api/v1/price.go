package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/service/tarot_service"
)

type SetPriceForm struct {
	SingleOrig       float32 `json:"single_orig" binding:"required"`                  // 单个原价
	SingleSellHigher float32 `json:"single_sell_higher" binding:"required"`           // 单个较高售价
	SingleSellLower  float32 `json:"single_sell_lower" binding:"required"`            // 单个较低售价
	ThreeOrig        float32 `json:"three_orig" binding:"required"`                   // 三个原价
	ThreeSellHigher  float32 `json:"three_sell_higher" binding:"required"`            // 三个较高售价
	ThreeSellLower   float32 `json:"three_sell_lower" binding:"required"`             // 三个较低售价
	Language         string  `json:"language" binding:"required" enums:"jp,zh,en,tc"` // 地区
}

type UpdatePriceForm struct {
	SingleOrig       float32 `json:"single_orig"`                                     // 单个原价
	SingleSellHigher float32 `json:"single_sell_higher"`                              // 单个较高售价
	SingleSellLower  float32 `json:"single_sell_lower"`                               // 单个较低售价
	ThreeOrig        float32 `json:"three_orig"`                                      // 三个原价
	ThreeSellHigher  float32 `json:"three_sell_higher"`                               // 三个较高售价
	ThreeSellLower   float32 `json:"three_sell_lower"`                                // 三个较低售价
	Language         string  `json:"language" enums:"jp,zh,en,tc" binding:"required"` // 地区
}

// SetPrice
// @Summary 设置价格
// @Param _ body SetPriceForm true "参数"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/tarot/price [post]
// @Security ApiKeyAuth
// @Tags Manager
func SetPrice(c *gin.Context) {
	var data SetPriceForm
	appG := app.Gin{C: c}
	if err := c.ShouldBindJSON(&data); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	priceData := map[string]interface{}{
		"language":           data.Language,
		"single_orig":        data.SingleOrig,
		"single_sell_higher": data.SingleSellHigher,
		"single_sell_lower":  data.SingleSellLower,
		"three_orig":         data.SingleOrig,
		"three_sell_higher":  data.ThreeSellHigher,
		"three_sell_lower":   data.ThreeSellLower,
		"created_by":         c.GetString("username"),
		"updated_by":         c.GetString("username"),
	}
	if err := tarot_service.SetPrice(priceData); err != nil {
		appG.Response(http.StatusOK, "设置价格失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// GetPrice
// @Summary 获取价格
// @Success 200 {object} []models.Price
// @Failure 500 {object} app.Response
// @Router /manager/tarot/price [get]
// @Security ApiKeyAuth
// @Tags Manager
func GetPrice(c *gin.Context) {
	appG := app.Gin{C: c}
	data, err := tarot_service.GetPriceTotal()
	if err != nil {
		appG.Response(http.StatusOK, "获取失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// UpdatePrice
// @Summary 修改价格
// @Param _ body UpdatePriceForm true "参数"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/tarot/price [put]
// @Security ApiKeyAuth
// @Tags Manager
func UpdatePrice(c *gin.Context) {
	var data UpdatePriceForm
	appG := app.Gin{C: c}
	if err := c.ShouldBindJSON(&data); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	priceData := map[string]interface{}{
		"language":           data.Language,
		"single_orig":        data.SingleOrig,
		"single_sell_higher": data.SingleSellHigher,
		"single_sell_lower":  data.SingleSellLower,
		"three_orig":         data.ThreeOrig,
		"three_sell_higher":  data.ThreeSellHigher,
		"three_sell_lower":   data.ThreeSellLower,
		"updated_by":         c.GetString("username"),
	}
	if err := tarot_service.UpdatePrice(priceData); err != nil {
		appG.Response(http.StatusOK, "设置价格失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
