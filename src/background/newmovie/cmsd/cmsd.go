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
	"background/newmovie/model"
	"background/newmovie/config"
	"background/common/constant"
	aapi "background/newmovie/controller/api"
	ccms "background/newmovie/controller/cms"
	"background/common/cache"
	"background/newmovie/service"

	"background/common/middleware"
	cmid "background/newmovie/middleware"

	_ "github.com/go-sql-driver/mysql"
	"background/shortvideo/setting"
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

	cacheRedisAddr, cacheRedisPwd, err := setting.GetCacheRedis(db)
	if err != nil {
		logger.Fatal(err)
		return
	}

	if err := cache.RedisTest(cacheRedisAddr, cacheRedisPwd); err != nil {
		logger.Fatal(err)
		return
	}

	dbMiddleware := middleware.GetDbPrepareHandler(config.GetDBName(), config.GetDBSource(), config.IsOrmLogEnabled())
	appVerifyMiddleware := cmid.AppVerifyHandler(model.AppTypeApp)
	//signatureMiddleware := middleware.SignatureVerifyHandler(false) // config.IsProductionEnv())

	service.InitCache(cacheRedisAddr, cacheRedisPwd)

	service.SetArea(config.GetAreaData())

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.OPTIONS("*f", func(c *gin.Context) {})

	cms := r.Group("cms")
	cms.POST("/upload", ccms.FileUpload)

	cms.Use(dbMiddleware,appVerifyMiddleware)
	{
		cms.POST("/install",aapi.InstallationHandler)
		cms.GET("/upgrade",aapi.UpgradeHandler)

		cms.GET("/video/list", aapi.VideoListHandler)
		cms.GET("/video", aapi.VideoDetailHandler)
		//cms.GET("/video/search", aapi.VideoSearchHandler)

		//cms.GET("/video/topsearch", aapi.VideoTopSearchHandler)

		cms.GET("/recommend", aapi.RecommendHandler)

		cms.GET("/page", aapi.PageHandler)

		cms.GET("/opinion", aapi.OpinionHandler)

		cms.POST("/digg", aapi.DiggHandler)
		cms.GET("/digglist", aapi.DiggListHandler)

		cms.GET("/guess", aapi.GuessListHandler)

		cms.POST("/user/stream/add", aapi.UserStreamAddHandler)
		cms.POST("/user/stream/update", aapi.UserStreamUpdateHandler)
		cms.POST("/user/stream/delete", aapi.UserStreamDeleteHandler)
		cms.GET("/user/stream/list", aapi.UserStreamListHandler)
		cms.POST("/user/want", aapi.UserWantHandler)

		cms.GET("/stream/list", aapi.StreamListHandler)
		cms.GET("/stream", aapi.StreamDetailHandler)

		cms.GET("/web", aapi.WebVideoHandler)

		cms.GET("/search", aapi.SearchHandler)
		cms.GET("/topsearch", aapi.TopSearchHandler)

		cms.GET("/notification", aapi.NotifcationHandler)
	}

	cms.Use(dbMiddleware)
	{
		r.GET("/login", ccms.AdminLoginHandler)
		cms.POST("/admin/login", ccms.AdminLoginHandler)

		cms.POST("/video/save", ccms.MovieSaveHandler)
		cms.POST("/script/save", ccms.ScriptSettingSaveHandler)
	}

	r.Static("/html", "/root/Git/e94/src/background/newmovie/html/")
	r.Static("/res", config.GetStaticRoot())
	r.Static("/doc", "/root/Git/e94/doc/newmovie/")

	//h := http.FileServer(http.Dir("/root/data/storage/movie/"))
	//http.Handle("/pic/", http.StripPrefix("/pic/", h)) // 启动静态文件服务
	//Header().Set("Expires", time.Now().Format("MON, 02 Jan 2006 15:04:05 GMT"))

	r.Run(":16882")

}


