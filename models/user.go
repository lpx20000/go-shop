package models

import (
	"encoding/json"
	"time"
)

type User struct {
	UserId         uint `gorm:"AUTO_INCREMENT"`
	OpenId         string
	NickName       string `gorm:"column:nickName"`
	AvatarUrl      string `gorm:"column:avatarUrl"`
	Gender         uint8
	Country        string
	Province       string
	City           string
	AddressId      uint
	WxappId        uint
	CreateTime     int64
	UpdateTime     int64
	UserAddress    []UserAddress `gorm:"foreignkey:UserId;association_foreignkey:UserId" json:"address,omitempty" `              //hasMany
	AddressDefault UserAddress   `gorm:"foreignkey:AddressId;association_foreignkey:AddressId" json:"addressDefault,omitempty" ` //belongsTo
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

func GetUserInfoByOpenId(token interface{}) (userInfo User) {
	Db.Where(map[string]interface{}{
		"open_id": token,
	}).
		Preload("UserAddress").
		Preload("AddressDefault").
		First(&userInfo)
	return
}

func Register(userInfo string, wxappId uint, openId string) (userId uint, err error) {
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
