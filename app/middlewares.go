package app

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// CheckToken check token
func CheckToken(ctx *gin.Context) {
	token := ctx.Request.Header.Get("X-Codeup-Token")
	if token != "io5eech5iewo1Ozo" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		ctx.Abort()
		return
	}
	ctx.Next()
}

// Recover 异常捕获
func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic: %v\n", r)
			debug.PrintStack()
			c.JSON(http.StatusOK, gin.H{
				"code": "1",
				"msg":  errorToString(r),
				"data": nil,
			})
			c.Abort()
		}
	}()
	c.Next()
}

func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}
