package app

import (
	"log"

	"github.com/gin-gonic/gin"

	"consul-webhook/pkg"
	"consul-webhook/response"
)

// CodeUp 处理CodeUp push tag事件
func CodeUp(c *gin.Context) {
	body, err := pkg.GetAliCode(c)
	if err != nil {
		response.NewResp().ErrorResponse(c, err.Error())
		return
	}

	tag := pkg.GetCurrentTag(body)
	log.Printf("tag: %v", tag)
	if pkg.IsDeleteTag(body) {
		log.Println("delete tag event, skip...")
		response.NewResp().JsonResponse(c, "success")
	}

	if err = RsyncConsul(tag); err != nil {
		response.NewResp().ErrorResponse(c, err.Error())
		return
	}
	response.NewResp().JsonResponse(c, "success")
}

// Rsync 手工调用此接口实现consul版本回退
func Rsync(c *gin.Context) {
	tag := c.Param("tag")
	log.Printf("tag: %v", tag)

	if err := RsyncConsul(tag); err != nil {
		response.NewResp().ErrorResponse(c, err.Error())
		return
	}
	response.NewResp().JsonResponse(c, "success")
}
