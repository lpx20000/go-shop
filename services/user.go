package services

import (
	"encoding/json"
	"shop/models"
	"shop/pkg/util"
	"time"

	"github.com/pkg/errors"
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

type registerUser struct {
	NickName  string `json:"nickName"`
	Gender    uint8  `json:"gender"`
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

var (
	url = "https://api.weixin.qq.com/sns/jscode2session"
)

func UserLogin(userInfo, code string, wxappId uint) (session string, userId uint, err error) {
	var (
		appInfo models.Wxapp
		openId  string
	)
	models.Db.Select("app_id, app_secret").First(&appInfo)

	if len(appInfo.AppId) == 0 || len(appInfo.AppSecret) == 0 {
		err = errors.New("请到 [后台-小程序设置] 填写appid 和 appsecret")
		return
	}
	if session, openId, err = getSessionFromWeiChat(code, appInfo.AppId, appInfo.AppSecret); err != nil {
		return
	}
	userId, err = register(userInfo, wxappId, openId)
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

func register(userInfo string, wxappId uint, openId string) (userId uint, err error) {
	var register registerUser

	if err = json.Unmarshal([]byte(userInfo), &register); err != nil {
		return
	}
	user := models.User{
		OpenId:     openId,
		NickName:   register.NickName,
		AvatarUrl:  register.AvatarUrl,
		Gender:     register.Gender,
		Country:    register.Country,
		Province:   register.Province,
		City:       register.City,
		AddressId:  0,
		WxappId:    wxappId,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
	err = models.Db.Create(&user).Error
	userId = user.UserId
	return
}
