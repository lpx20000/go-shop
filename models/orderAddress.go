package models

type OrderAddress struct {
	OrderAddressId uint       `gorm:"PRIMARY_KEY" json:"order_address_id"`
	Name           string     `json:"name"`
	Phone          string     `json:"phone"`
	ProvinceId     int        `json:"province_id"`
	CityId         int        `json:"city_id"`
	RegionId       int        `json:"region_id"`
	Detail         string     `json:"detail"`
	OrderId        int        `json:"-"`
	UserId         int        `json:"user_id"`
	WxappId        string     `json:"-"`
	RegionInfo     RegionInfo `json:"region"`
}

func (u *OrderAddress) AfterFind() error {
	u.RegionInfo, _ = GetRegionInfo(u.ProvinceId, u.RegionId, u.CityId)
	return nil
}
