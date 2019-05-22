package models

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

func GetUserInfoByOpenId(token interface{}) (userInfo User) {
	Db.Where(map[string]interface{}{
		"open_id": token,
	}).
		Preload("UserAddress").
		Preload("AddressDefault").
		First(&userInfo)
	return
}
