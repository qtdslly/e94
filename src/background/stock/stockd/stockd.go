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
	"background/stock/model"
	"background/stock/config"
	"background/common/constant"
	cc "background/stock/controller"

	_ "github.com/go-sql-driver/mysql"
	"background/common/middleware"
	"io/ioutil"
	"strings"
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

	model.InitModel(db)

	r := gin.New()

	gin.SetMode(gin.DebugMode)

	dbMiddleware := middleware.GetDbPrepareHandler(config.GetDBName(), config.GetDBSource(), config.IsOrmLogEnabled())

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.OPTIONS("*f", func(c *gin.Context) {})

	cms := r.Group("cms")
	cms.Use(dbMiddleware)
	{
		//load(config.GetTmplRoot() ,r)
		r.LoadHTMLGlob(config.GetTmplRoot() + "*.html")
		cms.GET("/stock/price", cc.StockPriceHandler)
		cms.GET("/stock/list", cc.StockListHandler)
	}
	r.Run(":16882")

}

func load(path string,r *gin.Engine){
	files, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Error(err)
		return
	}

	for _ ,file := range files{
		if file.Name() == "feeds" || file.Name() == "pic" || file.Name() == "js"{
			continue
		}
		if strings.Contains(file.Name(),".js"){
			continue
		}
		if file.IsDir(){
			load(path + "/" + file.Name(),r)
		}else{
			r.LoadHTMLGlob(path + "/" + file.Name())
		}

	}
}