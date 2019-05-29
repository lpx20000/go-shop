package services

import (
	"encoding/json"

	"shop/models"
	"shop/pkg/e"
	"shop/pkg/logging"
)

type Help struct {
	Base
	List []*models.WxappHelp `json:"list"`
}

func (h Help) getHelpKey() string {
	return e.CAHCHA_APP_HELP
}

func (h *Help) GetAppHelp() (err error) {
	var (
		key      string
		exist    bool
		dataByte []byte
	)
	key = h.getHelpKey()
	if dataByte, exist, err = h.getDataFromRedis(key); err != nil {
		logging.LogError(err)
		return
	}
	if exist {
		err = json.Unmarshal(dataByte, h)
		return
	}
	if h.List, err = models.GetHelp(); err != nil {
		logging.LogError(err)
		return
	}
	if err = h.setDataWithKey(key, h); err != nil {
		logging.LogError(err)
		return
	}
	return
}
