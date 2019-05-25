package models

type Wxapp struct {
	WxappId        uint        `gorm:"primary_key" json:"_"`
	AppName        string      `json:"_"`
	AppId          string      `json:"_"`
	AppSecret      string      `json:"_"`
	IsService      int         `json:"is_service"`
	ServiceImageId uint        `json:"_"`
	IsPhone        int         `json:"is_phone"`
	PhoneNo        string      `json:"phone_no"`
	PhoneImageId   uint        `json:"_"`
	Mchid          string      `json:"_"`
	Apikey         string      `json:"_"`
	ServiceImage   string      `json:"service_image"`
	PhoneImage     string      `json:"phone_image"`
	Navbar         WxappNavbar `gorm:"foreignkey:WxappId" json:"navbar"`
}

func (app *Wxapp) AfterFind() error {
	var upload []UploadFile
	Db.Model(&upload).Where("file_id in (?)", []uint{app.ServiceImageId, app.PhoneImageId}).Find(&upload)

	if len(upload) == 0 {
		return nil
	}
	for _, v := range upload {
		if app.ServiceImageId == v.FileId {
			app.ServiceImage = v.FileName
		}
		if app.PhoneImageId == v.FileId {
			app.PhoneImage = v.FileName
		}
	}

	return nil
}

func GetAppBase(appId uint) (app Wxapp, err error) {
	err = Db.Select("wxapp_id, is_service, service_image_id, is_phone, phone_no, phone_image_id").
		Where(&Wxapp{WxappId: appId}).
		Preload("Navbar").
		First(&app).Error
	return
}
