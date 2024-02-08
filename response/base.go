package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func (r *Resp) JsonResponse(c *gin.Context, data interface{}) {
	r.Code = http.StatusOK
	r.Msg = "success"
	r.Data = data
	c.JSON(r.Code, r)
}

func (r *Resp) ErrorResponse(c *gin.Context, data interface{}) {
	r.Code = http.StatusInternalServerError
	r.Msg = "failed"
	r.Data = data
	c.JSON(r.Code, r)
}

func NewResp() *Resp {
	return &Resp{}
}
