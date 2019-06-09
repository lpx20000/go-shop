package models

import (
	"fmt"
)

type UploadFile struct {
	FileId    uint   `json:"file_id,omitempty"`
	Storage   string `json:"storage,omitempty"`
	GroupId   uint   `json:"-"`
	FileUrl   string `json:"file_url"`
	FileName  string `json:"file_name,omitempty"`
	FileSize  string `json:"file_size,omitempty"`
	FileType  string `json:"file_type,omitempty"`
	Extension string `json:"extension,omitempty"`
	FilePath  string `json:"file_path,omitempty"`
	IsDelete  int    `json:"is_delete"`
	WxapId    uint   `json:"-,omitempty"`
}

func (upload *UploadFile) AfterFind() error {
	if upload.Storage == "local" {
		upload.FilePath = fmt.Sprintf("%s/%suploads/%s", Host, upload.FileUrl, upload.FileName)
	} else {
		upload.FilePath = fmt.Sprintf("%s/%s", upload.FileUrl, upload.FileName)
	}
	return nil
}
