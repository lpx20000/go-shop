package services

import "shop/models"

//获取分类
func GetCategory() (category []models.Category, err error) {
	category, err = models.GetCategoryInfo()
	return
}
