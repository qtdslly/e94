package main

import (
	"background/qrcode/controller"
	"background/qrcode/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	gin.SetMode(gin.DebugMode)

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CorsAllowHandler)
	r.OPTIONS("*f", func(c *gin.Context) {})

	baseApi := r.Group("")
	{
		baseApi.GET("/qrcode", controller.IptvQrcodeHandler)
	}

	r.Run(":6600") // listen and server on 0.0.0.0:9000
}
