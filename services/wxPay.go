package services

import "shop/models"

type WeChatPay struct {
	Mchid     string `json:"-"`
	AppId     string `json:"-"`
	AppSecret string `json:"-"`
	ApiKey    string `json:"-"`
}

func (w *WeChatPay) GetWxPayConfig(WxappId string) (err error) {
	var appInfo models.Wxapp
	if appInfo, err = models.GetAppPayInfo(WxappId); err != nil {
		return
	}
	w.AppId = appInfo.AppId
	w.Mchid = appInfo.Mchid
	w.ApiKey = appInfo.Apikey
	w.AppSecret = appInfo.AppSecret
	return
}

func (w *WeChatPay) WxPayNotify() (err error) {
	return nil
}
