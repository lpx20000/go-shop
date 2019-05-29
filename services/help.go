package services

import (
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/logging"
)

type Help struct {
	Base
	List []*models.WxappHelp `json:"list"`
}

func (h *Help) GetHelpKey() string {
	return e.CAHCHA_APP_HELP
}

func (h *Help) GetAppHelp() (err error) {
	var (
		key   string
		exist bool
	)
	key = h.GetHelpKey()
	if exist, err = h.GetDataFromRedis(key); err != nil {
		logging.LogError(err)
		return
	}
	if !exist {
		if h.List, err = models.GetHelp(); err != nil {
			logging.LogError(err)
			return
		}
		if err = h.SetDataWithKey(key, h); err != nil {
			logging.LogError(err)
			return
		}
	}
	return
}
