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
	"background/ams/model"
	"background/ams/config"
	"background/common/constant"
	//"background/common/cache"
	api "background/ams/controller/api"

	amsmid "background/ams/middleware"


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
	db.DB().SetMaxIdleConns(20)
	db.DB().SetMaxOpenConns(20)

	model.InitModel(db)

	r := gin.New()

	gin.SetMode(gin.DebugMode)

	dbMiddleware := middleware.GetDbPrepareHandler(config.GetDBName(), config.GetDBSource(), config.IsOrmLogEnabled())
	//signatureMiddleware := middleware.SignatureVerifyHandler(false) // config.IsProductionEnv())


	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.OPTIONS("*f", func(c *gin.Context) {})

	ams := r.Group("ams")

	r.GET("/", func(c *gin.Context) {
		url := config.GetDomain() + "/html/index.html"

		logger.Debug(url)
		c.Redirect(http.StatusPermanentRedirect, url)
		return
	})

	ams.Use(dbMiddleware)
	{
		ams.GET("/login", api.AdminLoginHandler)

	}

	ams.Use(dbMiddleware,amsmid.AdminVerifyHandler)
	{
		ams.GET("/admin/video/list", api.VideoListHandler)
		ams.POST("/admin/video/add", api.VideoAddHandler)
		ams.POST("/admin/video/delete", api.VideoDeleteHandler)

		ams.GET("/admin/episode/list", api.EpisodeListHandler)
		ams.POST("/admin/episode/add", api.EpisodeAddHandler)
		ams.POST("/admin/episode/delete", api.EpisodeDeleteHandler)

		ams.GET("/admin/playurl/list", api.PlayUrlListHandler)
		ams.POST("/admin/playurl/add", api.PlayUrlAddHandler)
		ams.GET("/admin/playurl/delete",api.PlayUrlDeleteHandler)
	}

	r.Static("/html", config.GetStaticRoot())

	r.Run(":8080")

}


