package services

import "shop/models"

func GetCategory() (category []models.Category, err error) {
	err = models.Db.Preload("Image").
		Order("sort ASC").
		Find(&category).Error
	return
}
