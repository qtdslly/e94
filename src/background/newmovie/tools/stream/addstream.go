package main

import (
	"background/newmovie/config"
	"background/common/logger"
	"flag"
	"background/newmovie/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	configPath := flag.String("conf", "../config/config.json", "Config file path")
	title := flag.String("t", "", "stream title")
	url := flag.String("u", "", "stream play url")
	category := flag.String("c", "", "caegory")

	flag.Parse()


	if len(*title) == 0 || len(*url) == 0{
		logger.Debug("useage: ./addstream -t 北京卫视 -u http://www.abc.m3u8")
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

	util.StreamAdd(*title,*url,*category,db)
}


