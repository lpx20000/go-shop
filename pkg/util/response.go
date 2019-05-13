package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop/pkg/e"
)

type R struct {
	Code int
	Data interface{}
}

func Response(c *gin.Context, r R) {
	c.JSON(http.StatusOK, gin.H{
		"code": r.Code,
		"msg":  e.GetMsg(r.Code),
		"data": r.Data,
	})
}
