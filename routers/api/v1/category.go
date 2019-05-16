package v1

import (
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/util"
	"sort"

	"github.com/gin-gonic/gin"
)

func GetGoodCategory(c *gin.Context) {
	var err error
	data := make(map[string]interface{})

	if info, err := models.GetCategory(); err == nil {
		childCategory := make(map[uint][]models.Category)
		parentCategory := make([]models.CategoryWithChild, 0)
		for _, category := range info {
			if category.ParentId != 0 {
				childCategory[category.ParentId] = append(childCategory[category.ParentId], category)
				continue
			}
			parentCategory = append(parentCategory, models.CategoryWithChild{
				Child:    make([]models.Category, 0),
				Category: category,
			})
		}

		for index, value := range parentCategory {
			if len(childCategory[value.CategoryId]) == 0 {
				continue
			}
			sort.Sort(models.Categories(childCategory[value.CategoryId]))
			parentCategory[index].Child = childCategory[value.CategoryId]
		}
		sort.Sort(models.CategoriesWithChild(parentCategory))
		data["list"] = parentCategory
		util.Response(c, util.R{Code: e.SUCCESS, Data: data})
		return
	}
	data["list"] = err
	util.Response(c, util.R{Code: e.ERROR, Data: data})
}

func Len() {

}
