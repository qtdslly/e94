package main

import (
	"os"
	"fmt"
	"flag"
	"html/template"
	"background/common/systemcall"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"background/common/logger"
	"log"
	"background/newmovie/model"
	"background/newmovie/config"
	"background/common/constant"
	aapi "background/newmovie/controller/api"
	ccms "background/newmovie/controller/cms"

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
		cms.POST("/install",aapi.InstallationHandler)

		cms.GET("/movie/list", aapi.NewMovieListHandler)
		cms.GET("/movie/search", aapi.NewMovieSearchHandler)
		cms.GET("/movie/topsearch", aapi.NewMovieTopSearchHandler)

		cms.POST("/admin/login", ccms.AdminLoginHandler)

		cms.POST("/movie/save", ccms.MovieSaveHandler)
	}

	h := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", h)) // 启动静态文件服务
	//Header().Set("Expires", time.Now().Format("MON, 02 Jan 2006 15:04:05 GMT"))
	http.HandleFunc("/", load)

	r.Run(":16882")

}


func load(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("*.html")
		if err != nil {
			log.Fatal(err)
		}

		data := struct {
			Title string
		}{
			Title: "品牌街-上天猫，就够了",
		}

		err = t.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
	} else {

	}
}

