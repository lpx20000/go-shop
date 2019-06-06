package util

import (
	"net/http"
	"shop/pkg/e"

	"github.com/gin-gonic/gin"
)

type R struct {
	Code int
	Msg  string
	Data interface{}
}

func Response(c *gin.Context, r R) {
	if len(r.Msg) == 0 {
		r.Msg = e.GetMsg(r.Code)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": r.Code,
		"msg":  r.Msg,
		"data": r.Data,
	})
}
