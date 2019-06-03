package services

import (
	"encoding/json"

	"shop/models"
	"shop/pkg/e"
	"shop/pkg/logging"
)

type AppBase struct {
	WxApp *models.Wxapp `json:"wxapp"`
}

func (b *AppBase) GetBaseKey() string {
	return e.CACHA_APP_BASE
}

func (b *AppBase) GetAppBase(appId uint) (err error) {
	var (
		key      string
		exist    bool
		dataByte []byte
	)
	key = b.GetBaseKey()
	if dataByte, exist, err = get(key); err != nil {
		logging.LogError(err)
		return
	}
	if exist {
		err = json.Unmarshal(dataByte, b)
		return
	}
	if b.WxApp, err = models.GetAppBase(appId); err != nil {
		logging.LogError(err)
		return
	}
	if err = set(key, b); err != nil {
		logging.LogError(err)
		return
	}
	return
}
