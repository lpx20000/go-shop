package models

type GoodsImage struct {
	Id         uint       `json:"id,omitempty"`
	GoodsId    uint       `json:"goods_id,omitempty"`
	ImageId    uint       `json:"image_id,omitempty"`
	WxappId    uint       `json:"-"`
	FilePath   string     `json:"file_path,omitempty"`
	FileName   string     `json:"file_name,omitempty"`
	FileUrl    string     `json:"file_url,omitempty"`
	UploadFile UploadFile `gorm:"foreignkey:FileId;association_foreignkey:ImageId" json:"-"`
}

func (g *GoodsImage) AfterFind() error {
	Db.Model(&g).Related(&g.UploadFile, "UploadFile")
	g.FilePath = g.UploadFile.FilePath
	g.FileName = g.UploadFile.FileName
	g.FileUrl = g.UploadFile.FileUrl
	return nil
}
