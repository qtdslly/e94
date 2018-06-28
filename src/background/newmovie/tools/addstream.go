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
	title := flag.String("title", "", "stream title")
	url := flag.String("url", "", "stream play url")
	provider := flag.Int("provider", "100", "provider")

	if len(*title) == 0 || len(*url) == 0 || *provider == 100{
		logger.Debug("useage: ./addstream -provider 8 -title 北京卫视 -url http://www.abc.m3u8")
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

	var stream model.Stream
	stream.Title = *title
	stream.Title = strings.Replace(stream.Title,"高清","",-1)
	stream.Title = strings.Replace(stream.Title,"-","",-1)
	stream.Title = TrimChinese(stream.Title)
	stream.Pinyin = util.TitleToPinyin(stream.Title)
	stream.Title = strings.Trim(stream.Title," ")
	logger.Debug(stream.Title)

	tx := db.Begin()
	if err := tx.Where("title = ?",stream.Title).First(&stream).Error ; err == gorm.ErrRecordNotFound{
		if strings.Contains(stream.Title,"CCTV"){
			stream.Category = "央视"
		}else if strings.Contains(stream.Title,"卫视"){
			stream.Category = "卫视"
		}else{
			stream.Category = "地方"
		}

		stream.OnLine = constant.MediaStatusOnLine
		stream.Sort = 0

		if err = tx.Create(&stream).Error ; err != nil{
			tx.Rollback()
			logger.Error(err)
			return
		}
	}

	var play model.PlayUrl
	play.Url = *url
	play.Provider = uint32(*provider)
	if err := tx.Where("provider = ? and url = ?",play.Provider,play.Url).First(&play).Error ; err == gorm.ErrRecordNotFound{
		play.Title = title
		play.OnLine = constant.MediaStatusOnLine
		play.Sort = 0
		play.ContentType = uint8(constant.MediaTypeStream)
		play.ContentId = stream.Id
		play.Quality = uint8(constant.VideoQuality720p)

		if err = tx.Create(&play).Error ; err != nil{
			tx.Rollback()
			logger.Error(err)
			return
		}
	}
}
