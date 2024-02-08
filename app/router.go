package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	prefix := "/api/v1/consul-webhook"
	api := r.Group(prefix, CheckToken)
	{
		api.POST("/codeup", CodeUp)
		api.POST("/rsync/tags/:tag", Rsync)
	}

	// health’s endpoint
	health := r.Group("/")
	{
		health.GET("/ping", func(context *gin.Context) {
			context.Writer.WriteString("pong")
		})
	}

	// 404
	r.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusNotFound, gin.H{
			"msg": "no endpoint",
		})
	})
}
