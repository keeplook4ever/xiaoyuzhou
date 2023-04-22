package util

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// GetPage get page parameters
func GetPage(c *gin.Context) int {
	result := 0
	page := com.StrTo(c.Query("page")).MustInt()
	pageSize := com.StrTo(c.Query("pageSize")).MustInt()
	if page > 0 {
		result = (page - 1) * pageSize
	}

	return result
}

// GetPageSize 获取每页展示量，默认10
func GetPageSize(c *gin.Context) int {
	result := 1000
	pageSize := com.StrTo(c.Query("pageSize")).MustInt()
	if pageSize > 0 {
		result = pageSize
	}
	return result
}
