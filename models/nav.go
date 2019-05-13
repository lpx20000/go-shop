package models

type WxappNavbar struct {
	WxappId            uint                   `gorm:"index" json:"_"`
	WxappTitle         string                 `json:"wxapp_title"`
	TopTextColor       int                    `json:"_"`
	TopTextColorBack   map[string]interface{} `json:"top_text_color"`
	TopBackgroundColor string                 `json:"top_background_color"`
}

func (nav *WxappNavbar) AfterFind() error {
	textColor := map[int]map[string]interface{}{
		10: {"text": "#000000", "value": 10},
		20: {"text": "#ffffff", "value": 20},
	}
	nav.TopTextColorBack = textColor[nav.TopTextColor]
	return nil
}
