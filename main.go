package main

import (
	"fmt"
	"log"
	"net/http"
	"shop/pkg/setting"
	"shop/routers"
)

func main() {
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatalln(s.ListenAndServe())
}
