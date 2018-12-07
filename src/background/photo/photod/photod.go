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
	"background/photo/model"
	"background/photo/config"
	"background/common/constant"
	ccms "background/photo/controller/cms"

	"background/common/middleware"

	_ "github.com/go-sql-driver/mysql"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

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

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.OPTIONS("*f", func(c *gin.Context) {})

	cms := r.Group("cms")
	cms.POST("/upload", ccms.FileUpload)

	cms.Use(dbMiddleware)
	{

	}

	cms.Use(dbMiddleware)
	{

	}

	logger.Debug(config.GetStaticRoot())
	r.Static("/html", config.GetStaticRoot())

	r.Run(":8000")

}


