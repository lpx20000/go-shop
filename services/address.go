package services

import (
	"encoding/json"
	"fmt"
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/logging"
	"strconv"
	"strings"
)

const (
	UPDATE = "update"
	CREATE = "create"
)

type DetailAddress struct {
	Detail models.UserAddress `json:"detail"`
	Region []string           `json:"region"`
}

type Address struct {
	List          []*models.UserAddress `json:"list"`
	DefaultId     int                   `json:"default_id"`
	AddressId     int                   `json:"-"`
	UserId        int                   `json:"-"`
	DetailAddress DetailAddress         `json:"-"`
	AddAddress    CommonAddress         `json:"-"`
	UpdateAddress UpdateAddress         `json:"-"`
}

type UpdateAddress struct {
	CommonAddress
	AddressId int `form:"address_id" binding:"required"  json:"address_id"`
}

type CommonAddress struct {
	UserId  int    `form:"user_id" json:"user_id"`
	WxappId string `form:"wxapp_id"  json:"wxapp_id"`
	Name    string `form:"name" binding:"required" json:"name"`
	Phone   string `form:"phone" binding:"required" json:"phone"`
	Detail  string `form:"detail" binding:"required" json:"detail"`
	Region  string `form:"region" binding:"required" json:"region"`
}

func (a *Address) getKey(userId int) string {
	return fmt.Sprintf("%s:%d", e.CACHA_APP_ADDRESS, userId)
}

func (a *Address) GetUserAddress(uid int) (err error) {
	var (
		dataByte []byte
		key      string
		exist    bool
		userInfo models.User
	)
	key = a.getKey(uid)

	if dataByte, exist, err = get(key); err != nil {
		logging.LogTrace(err)
		return
	}
	if userInfo, err = models.GetUserInfoByUid(uid); err != nil {
		logging.LogTrace(err)
		return
	}
	if exist {
		if err = json.Unmarshal(dataByte, a); err != nil {
			logging.LogTrace(err)
		}
		a.DefaultId = userInfo.AddressId
		return
	}
	if a.List, err = models.GetUserAddressList(uid); err != nil {
		logging.LogTrace(err)
		return
	}

	a.DefaultId = userInfo.AddressId
	if err = set(key, a); err != nil {
		logging.LogTrace(err)
	}
	return
}

func (a *Address) SetDefaultAddress() (err error) {
	err = models.UpdateUserInfo(a.UserId, map[string]interface{}{"address_id": a.AddressId})
	return
}

func (a *Address) DeleteAddress() (err error) {
	err = models.DeleteAddress(a.AddressId)
	return
}

func (a *Address) GetAddressDetail() (err error) {
	a.DetailAddress.Detail, err = models.GetAddressDetail(a.AddressId)
	if err != nil {
		logging.LogTrace(err)
		return
	}
	region := []string{
		a.DetailAddress.Detail.RegionInfo.Province,
		a.DetailAddress.Detail.RegionInfo.City,
		a.DetailAddress.Detail.RegionInfo.Region,
	}
	a.DetailAddress.Region = region
	return
}

func (a *Address) ModifyUserAddress(action string) (err error) {
	var (
		dataByte   []byte
		exist      bool
		key        string
		regionKey  string
		regionInfo map[string]int
	)
	regionKey = e.CACHA_APP_REGION_ID

	if dataByte, exist, err = get(regionKey); err != nil {
		logging.LogError(err)
		return
	}

	if exist {
		if err = json.Unmarshal(dataByte, &regionInfo); err != nil {
			logging.LogError(err)
			return
		}
	} else {
		if regionInfo, err = models.GetIdByRegionName(); err != nil {
			logging.LogError(err)
			return
		}
	}

	if action == CREATE {
		if err = a.createAddress(regionInfo); err != nil {
			logging.LogError(err)
			return
		}
	} else {
		if err = a.updateAddress(regionInfo); err != nil {
			logging.LogError(err)
			return
		}
	}

	key = a.getKey(a.UserId)
	if _, err = deleteCache(key); err != nil {
		logging.LogError(err)
	}
	return
}

func (a *Address) createAddress(regionInfo map[string]int) (err error) {
	var (
		userAddress models.UserAddress
		wxappId     int
		region      []string
	)
	wxappId, _ = strconv.Atoi(strings.TrimSpace(a.AddAddress.WxappId))
	region = strings.Split(strings.TrimSpace(a.AddAddress.Region), ",")
	userAddress.UserId = a.UserId
	userAddress.WxappId = wxappId
	userAddress.ProvinceId = regionInfo[fmt.Sprintf("%s:%d:%d", strings.TrimSpace(region[0]), 1, 0)]
	userAddress.CityId = regionInfo[fmt.Sprintf("%s:%d:%d", strings.TrimSpace(region[1]), 2, userAddress.ProvinceId)]
	userAddress.RegionId = regionInfo[fmt.Sprintf("%s:%d:%d", strings.TrimSpace(region[2]), 3, userAddress.CityId)]
	userAddress.Name = a.AddAddress.Name
	userAddress.Phone = a.AddAddress.Phone
	userAddress.Detail = a.AddAddress.Detail
	err = models.CreateUserAddress(userAddress)
	return
}

func (a *Address) updateAddress(regionInfo map[string]int) (err error) {
	var (
		userAddress models.UserAddress
	)
	wxappId, _ := strconv.Atoi(strings.TrimSpace(a.UpdateAddress.WxappId))
	region := strings.Split(strings.TrimSpace(a.UpdateAddress.Region), ",")
	userAddress.AddressId = a.UpdateAddress.AddressId
	userAddress.UserId = a.UpdateAddress.UserId
	userAddress.WxappId = wxappId
	userAddress.ProvinceId = regionInfo[fmt.Sprintf("%s:%d:%d", strings.TrimSpace(region[0]), 1, 0)]
	userAddress.CityId = regionInfo[fmt.Sprintf("%s:%d:%d", strings.TrimSpace(region[1]), 2, userAddress.ProvinceId)]
	userAddress.RegionId = regionInfo[fmt.Sprintf("%s:%d:%d", strings.TrimSpace(region[2]), 3, userAddress.CityId)]
	userAddress.Name = a.UpdateAddress.Name
	userAddress.Phone = a.UpdateAddress.Phone
	userAddress.Detail = a.UpdateAddress.Detail
	err = models.UpdateUserAddress(userAddress)
	return
}
