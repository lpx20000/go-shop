package services

import (
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/logging"
)

type AppBase struct {
	Base
	WxApp *models.Wxapp `json:"wxapp"`
}

func (b *AppBase) GetBaseKey() string {
	return e.CAHCHA_APP_BASE
}

func (b *AppBase) GetAppBase(appId uint) (err error) {
	var (
		key   string
		exist bool
	)
	key = b.GetBaseKey()
	if exist, err = b.GetDataFromRedis(key); err != nil {
		logging.LogError(err)
		return
	}
	if !exist {
		if b.WxApp, err = models.GetAppBase(appId); err != nil {
			logging.LogError(err)
			return
		}
		if err = b.SetDataWithKey(key, b); err != nil {
			logging.LogError(err)
			return
		}
	}
	return
}
