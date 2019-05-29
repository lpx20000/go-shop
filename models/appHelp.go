package models

import "github.com/jinzhu/gorm"

type WxappHelp struct {
	HelpId  uint   `json:"help_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Sort    uint8  `json:"sort"`
	WxappId uint   `json:"-"`
}

func GetHelp() ([]*WxappHelp, error) {
	var (
		appHelp []*WxappHelp
		err     error
	)
	err = Db.Order("sort ASC").Find(&appHelp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return appHelp, nil
}
