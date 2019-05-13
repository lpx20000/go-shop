package models

type UploadFile struct {
	FileId    uint   `json:"file_id"`
	Storage   string `json:"storage"`
	GroupId   uint   `json:"group_id"`
	FileUrl   string `json:"file_url"`
	FileName  string `json:"file_name"`
	FileSize  string `json:"file_size"`
	FileType  string `json:"file_type"`
	Extension string `json:"extension"`
	WxapId    uint   `json:"_"`
}
