package models

type WxappPage struct {
	PageId      uint        `json:"_"`
	PageType    int         `json:"_"`
	PageData    string      `json:"_"`
	NewPageData interface{} `json:"page_data"`
	WxapId      uint        `json:"_"`
}

type NewItems struct {
	Items interface{}
}
