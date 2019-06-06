package services

import (
	"encoding/json"
	"shop/models"
	"shop/pkg/util"

	"github.com/pkg/errors"
)

type User struct {
	UserId int    `json:"user_id"`
	Detail Detail `json:"-"`
}

type Detail struct {
	UserInfo   models.User    `json:"userInfo"`
	OrderCount map[string]int `json:"orderCount"`
}

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

func UserLogin(userInfo, code string, wxappId uint) (token string, userId int, err error) {
	var (
		appInfo AppInfo
		openId  string
	)
	models.Db.Model(&models.Wxapp{}).Select("app_id, app_secret").Scan(&appInfo)

	if len(appInfo.AppId) == 0 || len(appInfo.AppSecret) == 0 {
		err = errors.New("请到 [后台-小程序设置] 填写appid 和 appsecret")
		return
	}
	if openId, err = getSessionFromWeiChat(code, appInfo.AppId, appInfo.AppSecret); err != nil {
		return
	}

	userId = models.GetUserIdByToken(openId)
	if userId > 0 {
		token, err = util.GenerateToken(openId, userId)
	} else {
		if userId, err = models.Register(userInfo, wxappId, openId); err != nil {
			return
		}
		token, err = util.GenerateToken(openId, userId)
	}
	return
}

func (u *User) GetUserDetail() (err error) {
	u.Detail.UserInfo = models.GetUserInfoByOpenId(u.UserId)
	u.Detail.OrderCount = make(map[string]int, 2)
	u.Detail.OrderCount["payment"] = models.GetOrderCount(u.UserId, "payment")
	u.Detail.OrderCount["received"] = models.GetOrderCount(u.UserId, "received")
	return
}

func getSessionFromWeiChat(code, appId, appSecret string) (openid string, err error) {
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

	openid = sessionInfo.Openid
	return
}
