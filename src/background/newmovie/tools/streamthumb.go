package main

import (
	"background/newmovie/config"
	"background/common/logger"
	"background/newmovie/model"
	"background/common/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"flag"
	"fmt"
	"os/exec"
	"time"
	"os"
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
		time.Sleep(time.Minute * 5)
	}
}


func StreamThumb(db *gorm.DB){
	var err error
	var streams []model.Stream
	if err = db.Order("id asc").Find(&streams).Error ; err != nil{
		logger.Error(err)
		return
	}
	for _ , stream := range streams{
		flag := false
		flag1 := false
		var playUrls []model.PlayUrl
		if err = db.Order("sort asc").Where("content_type = 4 and content_id = ?",stream.Id).Find(&playUrls).Error ; err != nil{
			logger.Error(err)
			return
		}
		thumb := ""
		for _,playUrl := range playUrls{
			thumb = CheckStreamUrl(stream.Thumb,playUrl.Url)
			if thumb != ""{
				flag = true
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
			stream.Thumb = thumb
		}
		stream.OnLine = flag1
		if err = db.Save(&stream).Error ; err != nil{
			logger.Error(err)
			return
		}
	}
}
func CheckStreamUrl(sourceFileName,url string)string{
	c2 := make(chan string, 1)
	ffmpegAddr := "/usr/bin/ffmpeg"
	code := util.RandString(6)
	now := time.Now()
	fileName := fmt.Sprintf("%04d%02d%02d%02d%02d%02d",now.Year(),now.Month(),now.Day(),now.Hour(),now.Minute(),now.Second()) + code + ".jpg"

	go func() {
		cmdStr := fmt.Sprintf("%s -i '%s' -y -s 320x240 -vframes 1 /root/data/storage/stream/%s", ffmpegAddr, url,fileName)
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
			sourceFileName = strings.Replace(sourceFileName,"/res/stream/","",-1)
			err := os.Remove("/root/data/storage/stream/" + sourceFileName)
			if err != nil {
				logger.Error("file remove Error!",err)
			}
			return "/res/stream/" + fileName
		}else{
			return ""
		}
	case <-time.After(time.Second * 10):
		return ""
	}

	return ""
}


