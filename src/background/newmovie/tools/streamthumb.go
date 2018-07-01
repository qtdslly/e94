package main

import (
	"background/newmovie/config"
	"background/common/logger"
	"background/newmovie/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"flag"
	"fmt"
	"os/exec"
	"time"
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

	var streams []model.Stream
	if err = db.Order("id desc").Find(&streams).Error ; err != nil{
		logger.Error(err)
		return
	}
	for _ , stream := range streams{
		var playUrls []model.PlayUrl
		if err = db.Where("content_type = 4 and content_id = ?",stream.Id).Find(&playUrls).Error ; err != nil{
			logger.Error(err)
			return
		}
		for _,playUrl := range playUrls{
			if CheckStreamUrl(stream.Id,playUrl.Url){
				break
			}
		}
	}
}


func CheckStreamUrl(streamId uint32,url string)bool{
	c2 := make(chan string, 1)
	ffmpegAddr := "/usr/bin/ffmpeg"
	go func() {
		cmdStr := fmt.Sprintf("%s -i '%s' -y -s 320x240 -vframes 1 /data/www/dreamvideo/public/thumb/stream/%d.jpg", ffmpegAddr, url,streamId)
		fmt.Println(cmdStr)
		cmd := exec.Command("bash", "-c", cmdStr)

		if err := cmd.Run(); err == nil {
			c2 <- "success"
		}else{
			logger.Error(err)
			c2 <- "error"
		}
	}()
	select {
	case res := <-c2:
		if res == "success"{
			return true
		}else{
			return false
		}
	case <-time.After(time.Second * 10):
		return false
	}

	return false
}


