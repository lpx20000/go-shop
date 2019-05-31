package services

import (
	"encoding/json"

	"shop/models"
	"shop/pkg/e"
	"shop/pkg/logging"
)

type Help struct {
	List []*models.WxappHelp `json:"list"`
}

func (h Help) getHelpKey() string {
	return e.CACHA_APP_HELP
}

func (h *Help) GetAppHelp() (err error) {
	var (
		key      string
		exist    bool
		dataByte []byte
	)
	key = h.getHelpKey()
	if dataByte, exist, err = getDataFromRedis(key); err != nil {
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
	if err = setDataWithKey(key, h); err != nil {
		logging.LogError(err)
	}
	return
}
