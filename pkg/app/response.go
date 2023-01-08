package app

import (
	"github.com/gin-gonic/gin"

	"xiaoyuzhou/pkg/e"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode int, errCodeInfo interface{}, data interface{}) {
	var errInfo string
	var code int
	if value, ok := errCodeInfo.(int); ok {
		code = value
		errInfo = e.GetMsg(code)
	} else if value, ok := errCodeInfo.(string); ok {
		code = e.ERROR
		errInfo = value
	}
	g.C.JSON(httpCode, Response{
		Code: code,
		Msg:  errInfo,
		Data: data,
	})
	return
}
