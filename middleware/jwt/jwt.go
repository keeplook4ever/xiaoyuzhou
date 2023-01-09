package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"

	"github.com/gin-gonic/gin"

	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/util"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.GetHeader("X-Token")
		if token == "" {
			code = e.InvalidParams
		} else {
			clams, err := util.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = e.ErrorAuthCheckTokenTimeout
				default:
					code = e.ErrorAuthCheckTokenFail
				}
			} else {
				c.Set("username", clams.Username)
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}
		c.Next()
	}
}
