package services

import (
	"encoding/json"
	"github.com/pkg/errors"
	"shop/models"
	"shop/pkg/util"
)

type UserInfo struct {
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	Language  string `json:"language"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarUrl string `json:"avatarUrl"`
}

type sessionInfo struct {
	SessionKey string `json:"session_key"`
	Openid     string `json:"openid"`
}

type errorInfo struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type AppInfo struct {
	AppId     string
	AppSecret string
}

var (
	url = "https://api.weixin.qq.com/sns/jscode2session"
)

func UserLogin(userInfo, code string, wxappId uint) (session string, userId int, err error) {
	var (
		appInfo AppInfo
		openId  string
	)
	models.Db.Select("app_id, app_secret").Scan(&appInfo)

	if len(appInfo.AppId) == 0 || len(appInfo.AppSecret) == 0 {
		err = errors.New("请到 [后台-小程序设置] 填写appid 和 appsecret")
		return
	}
	if session, openId, err = getSessionFromWeiChat(code, appInfo.AppId, appInfo.AppSecret); err != nil {
		return
	}
	userId, err = models.Register(userInfo, wxappId, openId)
	return
}

func GetUserDetail(token string) (data map[string]interface{}) {
	var (
		userInfo   models.User
		orderCount map[string]int
	)
	data = make(map[string]interface{})
	orderCount = make(map[string]int)

	userInfo = models.GetUserInfoByOpenId(token)

	data["userInfo"] = userInfo
	orderCount["payment"] = models.GetCount(userInfo.UserId, "payment")
	orderCount["received"] = models.GetCount(userInfo.UserId, "received")
	data["orderCount"] = orderCount
	return
}

func getSessionFromWeiChat(code, appId, appSecret string) (session, openid string, err error) {
	var (
		result      []byte
		sessionInfo sessionInfo
		errMessage  errorInfo
	)

	requestInfo := map[string]string{
		"appid":      appId,
		"secret":     appSecret,
		"grant_type": "authorization_code",
		"js_code":    code,
	}

	if result, err = util.HttpGet(url, requestInfo); err != nil {
		return
	}

	if err = json.Unmarshal(result, &sessionInfo); err != nil {
		return
	}

	if len(sessionInfo.Openid) == 0 {
		if err = json.Unmarshal(result, &errMessage); err != nil {
			return
		}
		err = errors.New(errMessage.ErrMsg)
		return
	}

	session, err = util.GenerateToken(sessionInfo.Openid)
	openid = sessionInfo.Openid
	return
}
