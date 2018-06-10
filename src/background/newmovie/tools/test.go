package main

import (
	"background/newmovie/service"
	"background/newmovie/config"

	"background/common/logger"
	"flag"
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"

	"background/newmovie/model"
	"time"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	url := service.GetRealUrl("youku","http://v.youku.com/v_show/id_XMzU0ODk0MzQ0MA==.html",service.GetJsCode())

	jsCode := service.GetJsCode()

	configPath := flag.String("conf", "../config/config.json", "Config file path")
	err := config.LoadConfig(*configPath)
	if err != nil {
		return
	}

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}

	var set model.KvStore
	if err := db.Where("`key` = 'script_setting_key'").First(&set).Error ; err != nil{
		logger.Error(err)
		return
	}
	set.Value = jsCode
	now := time.Now()
	set.UpdatedAt = now
	if err = db.Save(&set).Error ; err != nil{
		logger.Error(err)
		return
	}

	logger.Debug(url)
}
