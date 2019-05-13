package models

type UserAddress struct {
	AddressId  uint   `json:"address_id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	ProvinceId uint   `json:"province_id"`
	CityId     uint   `json:"city_id"`
	RegionId   uint   `json:"region_id"`
	Detail     string `json:"detail"`
	UserId     uint   `json:"user_id"`
	WxappId    uint   `json:"wxapp_id"`
}
