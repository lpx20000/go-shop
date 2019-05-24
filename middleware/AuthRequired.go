package middleware

import (
	"shop/pkg/e"
	"shop/pkg/util"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		wxappId := c.Query("wxapp_id")

		if len(token) == 0 {
			util.Response(c, util.R{Code: e.ERROR_AUTH_TOKEN, Data: e.GetMsg(e.ERROR_AUTH_TOKEN)})
			c.Abort()
			return
		}
		session, err := util.ParseToken(token)
		if err != nil {
			util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
			c.Abort()
			return
		}
		c.Set(token, session.OpenId)
		c.Set("wxappId", wxappId)
		c.Next()
	}
}
