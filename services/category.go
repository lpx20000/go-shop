package services

import (
	"encoding/json"
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/logging"
	"sort"
)

type Category struct {
	List []*models.CategoryWithChild `json:"list"`
}

func (c *Category) getKey() string {
	return e.CACHA_APP_CATEGORY
}

func (c *Category) GetCategory() (err error) {
	var (
		dataByte []byte
		key      string
		exist    bool
	)
	key = c.getKey()
	if dataByte, exist, err = get(key); err != nil {
		logging.LogTrace(err)
		return
	}
	if exist {
		if err = json.Unmarshal(dataByte, c); err != nil {
			logging.LogTrace(err)
		}
		return
	}
	if err = c.GetCategoryList(); err != nil {
		logging.LogTrace(err)
		return
	}
	if err = set(key, c); err != nil {
		logging.LogTrace(err)
		return
	}
	return
}

func (c *Category) GetCategoryList() (err error) {
	var (
		categories []*models.Category
	)
	if categories, err = models.GetCategoryInfo(); err != nil {
		logging.LogTrace(err)
		return
	}
	childCategory := make(map[uint][]*models.Category)
	parentCategory := make([]*models.CategoryWithChild, 0)
	for _, category := range categories {
		if category.ParentId != 0 {
			childCategory[category.ParentId] = append(childCategory[category.ParentId], category)
			continue
		}
		parentCategory = append(parentCategory, &models.CategoryWithChild{
			Child:    make([]*models.Category, 0),
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
	c.List = parentCategory
	return
}
