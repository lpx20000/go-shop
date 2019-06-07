package services

import (
	"shop/models"
)

type WeChatPay struct {
	Mchid     string `json:"-"`
	AppId     string `json:"-"`
	AppSecret string `json:"-"`
	ApiKey    string `json:"-"`
}

func (w *WeChatPay) GetWxPayConfig(appId string) (err error) {
	var appInfo models.Wxapp
	if appInfo, err = models.GetAppPayInfo(appId); err != nil {
		return
	}
	w.AppId = appInfo.AppId
	w.Mchid = appInfo.Mchid
	w.ApiKey = appInfo.Apikey
	w.AppSecret = appInfo.AppSecret
	return
}

//func (w *WeChatPay) WxPayNotify() (err error) {
//	//appId, mchId, apiKey string, handler Handler, errorHandler ErrorHandler
//	mch := core.NewServer(w.AppId, w.Mchid, w.ApiKey, core.HandlerFunc(wxHandler),
//		core.ErrorHandlerFunc(wxErrorHandler))
//	return nil
//}
//
//func wxHandler(c *core.Context) {
//	log.Println(c)
//}
//
//func wxErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
//	fmt.Println("abc")
//}
