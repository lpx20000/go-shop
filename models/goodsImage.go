package models

type GoodsImage struct {
	Id         uint       `json:"id"`
	GoodsId    uint       `json:"goods_id"`
	ImageId    uint       `json:"image_id"`
	WxappId    uint       `json:"-"`
	FilePath   string     `json:"file_path"`
	FileName   string     `json:"file_name"`
	FileUrl    string     `json:"file_url"`
	UploadFile UploadFile `gorm:"foreignkey:FileId;association_foreignkey:ImageId" json:"-"`
}

func (g *GoodsImage) AfterFind() error {
	db.Model(&g).Related(&g.UploadFile, "UploadFile")
	g.FilePath = g.UploadFile.FilePath
	g.FileName = g.UploadFile.FileName
	g.FileUrl = g.UploadFile.FileUrl
	return nil
}
