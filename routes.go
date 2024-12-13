package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"version-2/controller"
)

func ListenAndServe() {
	app := gin.Default()

	app.GET("ping", controller.Ping)

	err := app.Run()
	if err != nil {
		log.Println(err)
		return
	}
}
