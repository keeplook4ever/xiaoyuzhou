package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"xiaoyuzhou/pkg/app"
)

// AddTrueWorld
// @Summary 添加真言
// @Param lang query string true "语言"
// @Param world query string true "真言"
// @Produce json
// @Accept json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/true-world [post]
// @Security ApiKeyAuth
// @Tags Manager
func AddTrueWorld(c *gin.Context) {
	appG := app.Gin{C: c}
	lang := c.Query("language")
	world := c.Query("world")
	fmt.Println(appG, lang, world)
}

func EditTrueWorld(c *gin.Context) {
	appG := app.Gin{C: c}
	fmt.Println(appG)
}

func GetTrueWorld(c *gin.Context) {
	appG := app.Gin{C: c}
	fmt.Println(appG)
}
