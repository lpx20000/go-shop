package services

import (
	"encoding/json"
	"shop/models"
)

func GetAppBase(appId uint) (app models.Wxapp, err error) {
	app, err = models.GetAppBase(appId)
	return
}

func GetAppIndex() (data map[string]interface{}) {
	data = make(map[string]interface{})
	data["items"] = getPageItem()
	data["newest"] = models.GetIndexNewestGood()
	data["best"] = models.GetIndexBestGoods()
	return
}

func getPageItem() interface{} {
	var item models.WxappPage
	models.Db.Select("page_data").First(&item)
	items := item.PageData
	var newItem models.NewItems
	if err := json.Unmarshal([]byte(items), &newItem); err != nil {
		return err.Error()
	}
	return newItem.Items
}
