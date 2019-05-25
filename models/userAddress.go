package models

type UserAddress struct {
	AddressId  uint       `json:"address_id"`
	Name       string     `json:"name"`
	Phone      string     `json:"phone"`
	ProvinceId int        `json:"province_id"`
	CityId     int        `json:"city_id"`
	RegionId   int        `json:"region_id"`
	Detail     string     `json:"detail"`
	UserId     uint       `json:"user_id"`
	WxappId    uint       `json:"-"`
	RegionInfo RegionInfo `json:"region"`
}

type RegionInfo struct {
	Province string `json:"province"`
	City     string `json:"city"`
	Region   string `json:"region"`
}

func (u *UserAddress) AfterFind() error {
	//var region []Region
	//Db.Model(&Region{}).Select("id, pid, name, level").Find(&region)
	//log.Println(region)
	//return nil
	return nil
}
