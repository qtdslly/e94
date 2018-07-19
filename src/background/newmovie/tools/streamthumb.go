package main

import (
	"background/newmovie/config"
	"background/common/logger"
	"background/newmovie/model"
	//"background/common/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"background/newmovie/util"
	"flag"
	"fmt"
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
	model.InitModel(db)

	for{
		StreamThumb(db)
	}
}


func StreamThumb(db *gorm.DB){
	var err error
	var streams []model.Stream
	if err = db.Order("category,title asc").Find(&streams).Error ; err != nil{
		logger.Error(err)
		return
	}
	for _ , stream := range streams{
		if stream.Id == 1698{
			continue
		}
		flag1 := false
		var playUrls []model.PlayUrl
		if err = db.Order("sort asc").Where("content_type = 4 and content_id = ?",stream.Id).Find(&playUrls).Error ; err != nil{
			logger.Error(err)
			return
		}
		thumb := ""
		for _,playUrl := range playUrls{
			thumb = util.CheckStream(playUrl.Url,config.GetStorageRoot() + "stream/" + fmt.Sprint(stream.Id) + ".jpg")
			if thumb != ""{
				playUrl.OnLine = true
				if err = db.Save(&playUrl).Error ; err != nil{
					logger.Error(err)
					return
				}
				flag1 = true
			}else{
				playUrl.OnLine = false
				if err = db.Save(&playUrl).Error ; err != nil{
					logger.Error(err)
					return
				}
			}
		}
		if thumb != ""{
			stream.Thumb = strings.Replace(thumb,config.GetStorageRoot(),"http://www.ezhantao.com/res/",-1)
		}
		stream.OnLine = flag1
		if err = db.Save(&stream).Error ; err != nil{
			logger.Error(err)
			return
		}
	}
}
