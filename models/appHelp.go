package models

import (
	"github.com/gin-gonic/gin"
	"strings"
	"github.com/jinzhu/gorm"
)

type WxappHelp struct {
	HelpId  uint   `json:"help_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Sort    uint8  `json:"sort"`
	WxappId uint   `json:"-"`
}

func GetHelp() ([]*WxappHelp, error) {
	var (
		appHelp []*WxappHelp
		err     error
	)
	err = Db.Order("sort ASC").Find(&appHelp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return appHelp, nil
}

func GetClientIP(c *gin.Context) string {
	ClientIP := c.ClientIP()
	//fmt.Println("ClientIP:", ClientIP)
	RemoteIP, _ := c.RemoteIP()
	//fmt.Println("RemoteIP:", RemoteIP)
	ip := c.Request.Header.Get("X-Forwarded-For")
	if strings.Contains(ip, "127.0.0.1") || ip == "" {
		ip = c.Request.Header.Get("X-real-ip")
	}
	if ip == "" {
		ip = "127.0.0.1"
	}
	if RemoteIP.String() != "127.0.0.1" {
		ip = RemoteIP.String()
	}
	if ClientIP != "127.0.0.1" {
		ip = ClientIP
	}
	return ip
}

