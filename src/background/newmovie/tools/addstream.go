package main

import (
	"background/newmovie/config"
	"background/common/logger"
	"background/newmovie/model"
	"background/common/util"

	"strings"
	"background/common/constant"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"flag"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	configPath := flag.String("conf", "../config/config.json", "Config file path")
	title := flag.String("t", "", "stream title")
	url := flag.String("u", "", "stream play url")
	category := flag.String("c", "", "caegory")

	flag.Parse()


	if len(*title) == 0 || len(*url) == 0{
		logger.Debug("useage: ./addstream -p 8 -t 北京卫视 -u http://www.abc.m3u8")
		return
	}
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

	tx := db.Begin()

	u := *url
	start := strings.Index(u,"m3u8")
	m3u8 := ""
	if start >= 0{
		m3u8 = u[start:start+4]
	}
	var sort uint32
	if m3u8 == "m3u8"{
		sort = uint32(1)
	}else{
		sort = uint32(10)
	}

	urlFound := false
	var play model.PlayUrl
	play.Url = *url
	if err := tx.Where("content_type = 4 and url = ?",play.Url).First(&play).Error ; err == gorm.ErrRecordNotFound{
	}else if err == nil{
		urlFound = true
	}

	sTitle := *title
	sTitle = strings.Replace(sTitle,"高清","",-1)
	sTitle = strings.Replace(sTitle,"-","",-1)
	sTitle = util.TrimChinese(sTitle)
	sTitle = strings.Trim(sTitle," ")
	if urlFound {
		var stream model.Stream
		if err := tx.Where("id = ?", play.ContentId).First(&stream).Error; err != nil {
			tx.Rollback()
			logger.Error(err)
			return
		}

		stream.Sort = sort
		stream.Title = sTitle
		stream.Pinyin = util.TitleToPinyin(stream.Title)
		logger.Debug(stream.Title)
		if err = tx.Save(&stream).Error ; err != nil{
			logger.Error(err)
			tx.Rollback()
			return
		}
		play.Sort = sort
		play.Title = *title
		if err = tx.Save(&play).Error ; err != nil{
			logger.Error(err)
			tx.Rollback()
			return
		}
	}else{
		var stream model.Stream
		stream.Title = sTitle
		if err := tx.Where("title = ?",stream.Title).First(&stream).Error ; err == gorm.ErrRecordNotFound{
			if strings.Contains(stream.Title,"CCTV"){
				stream.Category = "央视"
			}else if strings.Contains(stream.Title,"卫视"){
				stream.Category = "卫视"
			}else{
				stream.Category = *category
			}
			stream.Pinyin = util.TitleToPinyin(stream.Title)

			stream.OnLine = constant.MediaStatusOnLine
			stream.Sort = sort

			if err = tx.Create(&stream).Error ; err != nil{
				tx.Rollback()
				logger.Error(err)
				return
			}
		}

		play.Url = *url
		play.ContentId = stream.Id
		play.Provider = uint32(0)
		play.Title = *title
		play.OnLine = constant.MediaStatusOnLine
		play.Sort = sort
		play.ContentType = uint8(constant.MediaTypeStream)
		play.Quality = uint8(constant.VideoQuality720p)
		if err = tx.Create(&play).Error ; err != nil{
			tx.Rollback()
			logger.Error(err)
			return
		}
	}

	tx.Commit()
}


