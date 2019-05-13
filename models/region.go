package models

type Region struct {
	Id         uint   `json:"id"`
	Pid        uint   `json:"pid"`
	Shortname  string `json:"shortname"`
	Name       string `json:"name"`
	MergerName string `json:"merger_name"`
	Level      uint8  `json:"level"`
	Pinyin     string `json:"pinyin"`
	Code       string `json:"code"`
	ZipCode    string `json:"zip_code"`
	First      string `json:"first"`
	Lng        string `json:"lng"`
	Lat        string `json:"lat"`
}
