package models

type User struct {
	UserId     uint
	OpenId     string
	NickName   string `gorm:"column:nickName"`
	AvatarUrl  string `gorm:"column:avatarUrl"`
	Gender     uint8
	Country    string
	Province   string
	City       string
	AddressId  uint
	WxappId    uint
	CreateTime int64
	UpdateTime int64
}
