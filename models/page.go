package models

import "encoding/json"

type WxappPage struct {
	PageId      uint        `json:"_"`
	PageType    int         `json:"_"`
	PageData    string      `json:"_"`
	NewPageData interface{} `json:"page_data"`
	WxapId      uint        `json:"_"`
}

type NewItems struct {
	Items interface{} `json:"items"`
}

func GetPageItem() (data interface{}, err error) {
	var item WxappPage
	Db.Select("page_data").First(&item)
	items := item.PageData
	var newItem NewItems
	if err = json.Unmarshal([]byte(items), &newItem); err != nil {
		return
	}
	data = newItem.Items
	return
}
