package main

import (
	"fmt"
	"log"
	"net/http"
	"shop/models"
	"shop/pkg/gredis"
	"shop/pkg/logging"
	"shop/pkg/setting"
	"shop/routers"

	"github.com/gin-gonic/gin"
)

func init() {
	logging.Setup()
	setting.Setup()
	gredis.SetUp()
	models.Setup()
}

// @title Golang Gin API
// @version 1.0
// @description An example of gin
// @termsOfService https://github.com/EDDYCJY/go-gin-example
// @license.name MIT
// @license.url https://github.com/EDDYCJY/go-gin-example/blob/master/LICENSE
func main() {
	gin.SetMode(setting.ServerSetting.RunMode)
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routers.InitRouter(),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	logging.LogPanic(server.ListenAndServe())
}
