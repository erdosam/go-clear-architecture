package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"version-2/controller"
)

func ListenAndServe() {
	app := gin.Default()

	app.HEAD("/ping", controller.PingHead)
	app.GET("/ping", controller.Ping)

	authorized := app.Group("/", gin.BasicAuth(gin.Accounts{
		"foo": "bar",
	}))

	// User's cart manager APIs
	cart := authorized.Group("/cart")
	{
		cart.GET("", controller.CartList)
	}

	err := app.Run()
	if err != nil {
		log.Println(err)
		return
	}
}
