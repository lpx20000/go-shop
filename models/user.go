package models

type User struct {
	UserId    uint   `json:"user_id"`
	OpenId    string `json:"open_id"`
	NickName  string `json:"nick_name"`
	AvatarUrl string `json:"avatarUrl"`
	Gender    uint8  `json:"gender"`
	Country   string `json:"country"`
	Province  string `json:"province"`
	City      string `json:"city"`
	AddressId uint   `json:"address_id"`
	WxappId   uint   `json:"wxapp_id"`
}
