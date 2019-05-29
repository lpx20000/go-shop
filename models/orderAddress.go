package models

type OrderAddress struct {
	OrderAddressId uint   `json:"order_address_id"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	ProvinceId     int    `json:"province_id"`
	CityId         int    `json:"city_id"`
	RegionId       int    `json:"region_id"`
	Detail         string `json:"detail"`
	OrderId        uint   `json:"-"`
	UserId         int    `json:"user_id"`
	WxappId        uint   `json:"-"`
}
