package services

import (
	"encoding/json"
	"shop/models"
)

func GetAppBase(appId uint) (app models.Wxapp, err error) {
	err = models.Db.Select("wxapp_id, is_service, service_image_id, is_phone, phone_no, phone_image_id").
		Where(&models.Wxapp{WxappId: appId}).
		Preload("Navbar").
		First(&app).Error
	return
}

func GetPageItem() interface{} {
	var item models.WxappPage
	models.Db.Select("page_data").First(&item)
	items := item.PageData
	var newItem models.NewItems
	if err := json.Unmarshal([]byte(items), &newItem); err != nil {
		return err.Error()
	}
	return newItem.Items
}

func GetNewestGood() (goods models.Goods) {
	models.Db.Where(map[string]interface{}{"is_delete": 0, "goods_status": models.ON_SALES}).
		Preload("Category").
		Preload("GoodsSpec").
		Preload("GoodsImage").
		Order("goods_id DESC").
		Order("goods_sort ASC").
		First(&goods)
	return
}

func GetBestGoods() (goods []models.Goods) {
	models.Db.Where(map[string]interface{}{"is_delete": 0, "goods_status": models.ON_SALES}).
		Preload("Category").
		Preload("GoodsSpec").
		Preload("GoodsImage").
		Order("goods_id DESC").
		Order("goods_sort ASC").
		Limit(10).
		Find(&goods)
	return
}
