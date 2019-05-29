package models

type UserAddress struct {
	AddressId  int        `gorm:"primary_key" json:"address_id"`
	Name       string     `json:"name"`
	Phone      string     `json:"phone"`
	ProvinceId int        `json:"province_id"`
	CityId     int        `json:"city_id"`
	RegionId   int        `json:"region_id"`
	Detail     string     `json:"detail"`
	UserId     int        `json:"user_id"`
	WxappId    int        `json:"-"`
	RegionInfo RegionInfo `json:"region"`
	Model
}

func (u *UserAddress) AfterFind() error {
	u.RegionInfo = GetRegionInfo(u.ProvinceId, u.CityId, u.RegionId)
	return nil
}

func GetUserAddressList(uid int) (address []UserAddress) {
	Db.Where(&UserAddress{UserId: uid}).Order("address_id ASC").Find(&address)
	return
}

func CreateUserAddress(userAddress UserAddress) (err error) {
	err = Db.Create(&userAddress).Error
	return
}