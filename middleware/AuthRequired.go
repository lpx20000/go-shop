package middleware

import (
	"shop/pkg/e"
	"shop/pkg/util"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token   string
			wxappId string
			err     error
		)

		if c.Request.Method == "GET" {
			token = strings.TrimSpace(c.Query("token"))
			wxappId = strings.TrimSpace(c.Query("wxapp_id"))
		} else if c.Request.Method == "POST" {
			token = strings.TrimSpace(c.Request.FormValue("token"))
			wxappId = strings.TrimSpace(c.Request.FormValue("wxapp_id"))
		}

		if len(token) == 0 || len(wxappId) == 0 {
			util.Response(c, util.R{Code: e.ERROR, Data: e.GetMsg(e.ERROR_AUTH_TOKEN)})
			c.Abort()
			return
		}
		claims, err := util.ParseToken(token)
		if err != nil {
			util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
			c.Abort()
			return
		}

		c.Set(token, claims.OpenId)
		c.Set("token", claims.OpenId)
		c.Set("userId", claims.UserId)
		c.Set("wxappId", wxappId)
		c.Next()
	}
}
