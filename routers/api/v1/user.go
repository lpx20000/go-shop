package v1

import (
	"shop/pkg/e"
	"shop/pkg/util"
	"shop/services"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
)

type User struct {
	WxappId       uint   `form:"wxapp_id" binding:"required" json:"wxapp_id"`
	Token         string `form:"token" json:"token"`
	Code          string `form:"code" binding:"required" json:"code"`
	UserInfo      string `form:"user_info" binding:"required" json:"user_info"`
	EncryptedData string `form:"encrypted_data" binding:"required" json:"encrypted_data"`
	Signature     string `form:"signature" binding:"required" json:"signature"`
	Iv            string `form:"iv" binding:"required" json:"iv"`
}

type authToken struct {
	Token string `form:"token" binding:"required"`
}

func UserLogin(c *gin.Context) {
	var (
		user   User
		userId int
		token  string
		err    error
	)

	//tokens, _ := util.GenerateToken("oZ54Q5T2Nzdhz6kLdfSSFwHBprwg")
	//util.Response(c, util.R{Code: e.INVALID_PARAMS, Data: tokens})
	//return

	data := make(map[string]interface{})
	if err = c.ShouldBindWith(&user, binding.FormPost); err != nil {
		util.Response(c, util.R{Code: e.INVALID_PARAMS, Data: err.Error()})
		return
	}

	token, userId, err = services.UserLogin(user.UserInfo, user.Code, user.WxappId)

	if err != nil {
		util.Response(c, util.R{Code: e.INVALID_PARAMS, Data: err.Error()})
		return
	}

	data["token"] = token
	data["user_id"] = userId
	util.Response(c, util.R{Code: e.SUCCESS, Data: data})
}

func GetUserDetail(c *gin.Context) {
	var (
		token string
		auth  authToken
		data  map[string]interface{}
	)
	if c.ShouldBindQuery(&auth) != nil {
		util.Response(c, util.R{Code: e.INVALID_PARAMS, Data: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	token = c.GetString(auth.Token)
	if token == "" {
		util.Response(c, util.R{Code: e.INVALID_PARAMS, Data: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	data = services.GetUserDetail(token)
	util.Response(c, util.R{Code: e.SUCCESS, Data: data})
}
