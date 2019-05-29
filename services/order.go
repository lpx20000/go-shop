package services

import "shop/models"

func GetOrderList(userId int, filter string) (data map[string]interface{}) {
	data = make(map[string]interface{})
	data["list"] = models.GetOrderList(userId, filter)
	return
}
