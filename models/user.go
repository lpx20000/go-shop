package models

import (
	"encoding/json"
	"time"
)

type User struct {
	UserId         int    `gorm:"AUTO_INCREMENT"`
	OpenId         string `json:"-"`
	NickName       string `gorm:"column:nickName"`
	AvatarUrl      string `gorm:"column:avatarUrl"`
	Gender         uint8
	Country        string
	Province       string
	City           string
	AddressId      uint
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

func GetUserInfoByOpenId(token interface{}) (userInfo User) {
	Db.Where(map[string]interface{}{
		"open_id": token,
	}).
		Preload("UserAddress").
		Preload("AddressDefault").
		First(&userInfo)
	return
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

func (u *User) AfterFind() error {
	var (
		all        map[int]CommonList
		commonList []CommonList
	)
	if u.AddressDefault.UserId > 0 {
		all = make(map[int]CommonList)

		commonList = GetRegion()

		for _, item := range commonList {
			all[item.Id] = item
		}
		u.AddressDefault.RegionInfo = RegionInfo{
			Province: all[u.AddressDefault.ProvinceId].Name,
			City:     all[u.AddressDefault.CityId].Name,
			Region:   all[u.AddressDefault.RegionId].Name,
		}
	}

	if len(u.UserAddress) > 0 {
		for index, address := range u.UserAddress {
			u.UserAddress[index].RegionInfo = RegionInfo{
				Province: all[address.ProvinceId].Name,
				City:     all[address.CityId].Name,
				Region:   all[address.RegionId].Name,
			}
		}
	}

	return nil
}

func GetUserIdByToken(token string) (uid int) {
	var user User
	Db.Where(&User{OpenId: token}).Select("user_id").First(&user)
	uid = user.UserId
	return
}
