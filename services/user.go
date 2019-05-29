package services

import (
	"encoding/json"
	"shop/models"
	"shop/pkg/util"
	"strconv"
	"strings"

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

type AddAddress struct {
	UserId int     `form:"user_id" json:"user_id"`
	WxappId string    `form:"wxapp_id"  json:"wxapp_id"`
	Name    string `form:"name" binding:"required" json:"name"`
	Phone   string `form:"phone" binding:"required" json:"phone"`
	Detail  string `form:"detail" binding:"required" json:"detail"`
	Region  string `form:"region" binding:"required" json:"region"`
}


var (
	url = "https://api.weixin.qq.com/sns/jscode2session"
)

func UserLogin(userInfo, code string, wxappId uint) (session string, userId int, err error) {
	var (
		appInfo AppInfo
		openId  string
	)
	models.Db.Model(&models.Wxapp{}).Select("app_id, app_secret").Scan(&appInfo)

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
	orderCount["payment"] = models.GetOrderCount(userInfo.UserId, "payment")
	orderCount["received"] = models.GetOrderCount(userInfo.UserId, "received")
	data["orderCount"] = orderCount
	return
}

func GetUserAddress(uid int) (data map[string]interface{}) {
	data = make(map[string]interface{})
	data["list"] = models.GetUserAddressList(uid)
	data["default_id"] = models.GetUserInfoByUid(uid).AddressId
	return
}

func AddUserAddress(address AddAddress) (err error)  {
	var (
		region []string
		provinceId int
		cityId int
		regionId int
		userAddress models.UserAddress
		wxappId int
	)
	wxappId, _ = strconv.Atoi(strings.TrimSpace(address.WxappId))
	region = strings.Split(strings.TrimSpace(address.Region), ",")
	provinceId = models.GetIdByRegionName(strings.TrimSpace(region[0]), 1, 0)
	cityId = models.GetIdByRegionName(strings.TrimSpace(region[1]), 2, provinceId)
	regionId = models.GetIdByRegionName(strings.TrimSpace(region[2]), 3, cityId)
	userAddress.UserId = address.UserId
	userAddress.WxappId = wxappId
	userAddress.ProvinceId = provinceId
	userAddress.CityId = cityId
	userAddress.RegionId = regionId
	userAddress.Name = address.Name
	userAddress.Phone = address.Phone
	userAddress.Detail = address.Detail
	err = models.CreateUserAddress(userAddress)
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
