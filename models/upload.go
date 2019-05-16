package models

import "fmt"

type UploadFile struct {
	FileId    uint   `json:"file_id"`
	Storage   string `json:"storage"`
	GroupId   uint   `json:"group_id"`
	FileUrl   string `json:"file_url"`
	FileName  string `json:"file_name"`
	FileSize  string `json:"file_size"`
	FileType  string `json:"file_type"`
	Extension string `json:"extension"`
	FilePath  string `json:"file_path"`
	WxapId    uint   `json:"-"`
}

func (upload *UploadFile) AfterFind() error {
	if upload.Storage == "local" {
		upload.FilePath = fmt.Sprintf("%s/%suploads/%s", HostInfo.Host, upload.FileUrl, upload.FileName)
	} else {
		upload.FilePath = fmt.Sprintf("%suploads/%s", upload.FileUrl, upload.FileName)
	}
	return nil
}
