package services

import (
	"encoding/json"

	"shop/models"
	"shop/pkg/e"
)

type App struct {
	Items  interface{}     `json:"items"`
	Newest []*models.Goods `json:"newest"`
	Best   []*models.Goods `json:"best"`
}

func (a *App) GetItemKey() string {
	return e.CACHE_APP_ITEM
}

func (a *App) getIndexKey() string {
	return e.CACHE_APP_INDEX
}

func (a *App) GetIndexData() (err error) {
	var (
		key      string
		exist    bool
		dataByte []byte
	)
	key = a.getIndexKey()
	if dataByte, exist, err = get(key); err != nil {
		return
	}

	if exist {
		err = json.Unmarshal(dataByte, a)
		return
	}
	if a.Items, err = models.GetPageItem(); err != nil {
		return
	}
	if a.Newest, err = models.GetIndexNewestGood(); err != nil {
		return
	}
	if a.Best, err = models.GetIndexBestGoods(); err != nil {
		return
	}
	err = set(key, a)
	return
}
