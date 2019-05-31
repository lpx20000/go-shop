package models

import "github.com/jinzhu/gorm"

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
	u.RegionInfo, _ = GetRegionInfo(u.ProvinceId, u.RegionId, u.CityId)
	return nil
}

func GetUserAddressList(uid int) ([]*UserAddress, error) {
	var (
		address []*UserAddress
		err     error
	)
	err = Db.Where(&UserAddress{UserId: uid}).Order("address_id DESC").Find(&address).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return address, nil
}

func CreateUserAddress(userAddress UserAddress) error {
	return Db.Model(&userAddress).Create(&userAddress).Error
}

func UpdateUserAddress(userAddress UserAddress) error {
	return Db.Model(&userAddress).Update(&userAddress).Error
}

func GetAddressDetail(addressId int) (UserAddress, error) {
	var (
		address UserAddress
		err     error
	)
	err = Db.Where(&UserAddress{AddressId: addressId}).First(&address).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return address, err
	}
	return address, nil
}

func DeleteAddress(addressId int) error {
	return Db.Delete(&UserAddress{AddressId: addressId}).Error
}
