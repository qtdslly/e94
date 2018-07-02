package main

import (
	"background/newmovie/config"
	"background/common/logger"
	"background/newmovie/model"
	"background/common/util"

	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"flag"
	"strings"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	configPath := flag.String("conf", "../config/config.json", "Config file path")
	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		return
	}

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}
	db.LogMode(true)

	var streams []model.Stream
	if err = db.Find(&streams).Error ; err != nil{
		logger.Error(err)
		return
	}

	for _, stream := range streams{
		var playUrls []model.StreamSource
		if err = db.Where("stream_id = ? ",stream.Id).Find(&playUrls).Error ; err != nil{
			logger.Error(err)
			continue
		}

		for _ , playUrl := range playUrls{
			if strings.Contains(playUrl.Url,"migu") || strings.Contains(playUrl.Url,"starschinalive"){
				continue
			}
			if !strings.Contains(playUrl.Url,"m3u8") && !strings.Contains(playUrl.Url,"rtmp"){
				continue
			}
			if util.CheckStreamUrl(playUrl.Url){
				fmt.Println("SUCCESS" + stream.Title + "\t" + playUrl.Url)
			}
		}
	}


}


