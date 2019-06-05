package main

import (
	"os"
	"fmt"
	"flag"
	"background/common/systemcall"
	"github.com/gin-gonic/gin"
	"background/common/logger"
	"log"
	"background/wdkq/config"
	"background/common/constant"

	"background/wdkq/api"

	"net/http"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	//check version
	if len(os.Args) > 1 {
		if os.Args[1] == "-version" {
			fmt.Println(constant.Version)
			os.Exit(0)
		} else if os.Args[1] != "-conf" {
			fmt.Println("invalid argument, only -conf/-version are accepted!")
			os.Exit(0)
		}
	}

	configPath := flag.String("conf", "./config/config.json", "Config file path")

	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("Config Failed!!!!", err)
		return
	}

	// set file descriptor limit
	systemcall.SetFileLimit()

	r := gin.New()

	gin.SetMode(gin.DebugMode)


	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.OPTIONS("*f", func(c *gin.Context) {})


	r.GET("/", func(c *gin.Context) {
		url := config.GetDomain() + "/html/index.html"

		logger.Debug(url)
		c.Redirect(http.StatusPermanentRedirect, url)
		return
	})

	r.Use()
	{
		r.POST("/save", api.SavePrintHandler)
	}

	r.Static("/html", config.GetStaticRoot())

	r.Run(":80")

}


