package main

import (
	"os"
	"fmt"
	"flag"
	"background/common/systemcall"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"background/common/logger"
	"log"
	"background/ezhantao/config"
	"background/common/constant"
  taobao "background/ezhantao/controller/api/taobao"
  wechart "background/ezhantao/controller/api/wechart"

	"background/common/middleware"

	_ "github.com/go-sql-driver/mysql"
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

	configPath := flag.String("conf", "../config/config.json", "Config file path")
	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("Config Failed!!!!", err)
		return
	}

	// set file descriptor limit
	systemcall.SetFileLimit()

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}
	db.DB().SetMaxIdleConns(10)

	r := gin.New()

	gin.SetMode(gin.DebugMode)

	dbMiddleware := middleware.GetDbPrepareHandler(config.GetDBName(), config.GetDBSource(), config.IsOrmLogEnabled())

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
		r.GET("/wechart", wechart.CheckHandler)
		r.POST("/wechart", wechart.WeChartHandler)
		r.GET("/play", wechart.PlayHandler)

	}
	cms := r.Group("cms")

	cms.Use(dbMiddleware)
	{
    cms.GET("/goods/list",taobao.GoodsListHandler)
    cms.GET("/tpwd",taobao.TPwdHandler)
    cms.GET("/search",taobao.SearchHandler)

  }
	r.Static("/html", config.GetStaticRoot())

	r.Run(":80")


}


