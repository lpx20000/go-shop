package models

type UploadGroup struct {
	GroupId   uint   `json:"group_id"`
	GroupType string `json:"group_type"`
	GroupName string `json:"group_name"`
	Sort      uint8  `json:"sort"`
	WxappId   uint   `json:"wxapp_id"`
}
