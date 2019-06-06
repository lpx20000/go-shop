package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	UserId         int    `gorm:"AUTO_INCREMENT,PRIMARY_KEY" json:"user_id"`
	OpenId         string `json:"-"`
	NickName       string `gorm:"column:nickName"`
	AvatarUrl      string `gorm:"column:avatarUrl"`
	Gender         interface{}
	Country        string        `json:"country"`
	Province       string        `json:"province"`
	City           string        `json:"city"`
	AddressId      int           `json:"address_id"`
	WxappId        uint          `json:"-"`
	CreateTime     int64         `json:"-"`
	UpdateTime     int64         `json:"-"`
	UserAddress    []UserAddress `gorm:"foreignkey:UserId;association_foreignkey:UserId" json:"address,omitempty" `               //hasMany
	AddressDefault UserAddress   `gorm:"foreignkey:AddressId;association_foreignkey:AddressId" json:"address_default,omitempty" ` //belongsTo
}

type RegisterUser struct {
	NickName  string `json:"nickName"`
	Gender    uint8  `json:"gender"`
	Language  string `json:"language"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarUrl string `json:"avatarUrl"`
}

func GetUserInfoByOpenId(uid int) User {
	var userInfo User
	Db.Where(&User{UserId: uid}).
		Preload("UserAddress").
		Preload("AddressDefault").
		First(&userInfo)
	return userInfo
}

func Register(userInfo string, wxappId uint, openId string) (userId int, err error) {
	var register RegisterUser

	if err = json.Unmarshal([]byte(userInfo), &register); err != nil {
		return
	}
	user := User{
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
	err = Db.Create(&user).Error
	userId = user.UserId
	return
}

func (u *User) AfterFind() (err error) {
	if u.AddressDefault.UserId > 0 {
		u.AddressDefault.RegionInfo, err = GetRegionInfo(u.AddressDefault.ProvinceId, u.AddressDefault.RegionId, u.AddressDefault.CityId)
		if err != nil {
			return
		}
		if u.Gender == 1 {
			u.Gender = "男"
		} else {
			u.Gender = "女"
		}
	}

	if len(u.UserAddress) > 0 {
		for index, address := range u.UserAddress {
			u.UserAddress[index].RegionInfo, err = GetRegionInfo(address.ProvinceId, address.RegionId, address.CityId)
			if err != nil {
				return
			}
		}
	}

	return nil
}

func GetUserIdByToken(openId string) (uid int) {
	var user User
	Db.Where(&User{OpenId: openId}).Select("user_id").First(&user)
	uid = user.UserId
	return
}

func GetUserInfoByUid(uid int) (userInfo User, err error) {
	err = Db.Where(&User{UserId: uid}).First(&userInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return
}

func UpdateUserInfo(uid int, updates map[string]interface{}) (err error) {
	err = Db.Model(&User{}).Where(&User{UserId: uid}).Updates(updates).Error
	return
}
