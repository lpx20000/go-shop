package models

type AppHelp struct {
	HelpId  uint   `json:"help_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Sort    uint8  `json:"sort"`
	WxappId uint   `json:"wxapp_id"`
}
