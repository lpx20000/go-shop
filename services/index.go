package services

import (
	"shop/models"
	"shop/pkg/e"
)

type App struct {
	Base
	Items  interface{}     `json:"items"`
	Newest []*models.Goods `json:"newest"`
	Best   []*models.Goods `json:"best"`
}

func (a *App) GetItemKey() string {
	return e.CACHE_APP_ITEM
}

func (a *App) GetIndexKey() string {
	return e.CAHCHE_APP_INDEX
}

func (a *App) GetIndexData() (err error) {
	var (
		key   string
		exist bool
	)
	key = a.GetIndexKey()
	if exist, err = a.GetDataFromRedis(key); err != nil {
		return
	}

	if !exist {
		if a.Items, err = models.GetPageItem(); err != nil {
			return
		}
		if a.Newest, err = models.GetIndexNewestGood(); err != nil {
			return
		}
		if a.Best, err = models.GetIndexBestGoods(); err != nil {
			return
		}
		err = a.SetDataWithKey(key, a)
	}
	return
}
