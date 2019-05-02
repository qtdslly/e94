package main

import (
	"background/common/util"
	"background/common/logger"
	"background/newmovie/config"
	"background/newmovie/model"

	"time"
	"errors"
	"flag"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)
func main(){
	logger.SetLevel(config.GetLoggerLevel())
	configPath := flag.String("conf", "../config/config.json", "Config file path")
	flag.Parse()


	err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Error(err)
		return
	}

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}
	db.LogMode(true)
	model.InitModel(db)
	var playUrls []model.PlayUrl
	if err := db.Where("content_type = 4").Find(&playUrls).Error ; err != nil{
		logger.Error(err)
		return
	}

	for _,playUrl := range playUrls{
		width,height,ready,err := GetStreamRation(playUrl.Url)
		if err == nil{
			playUrl.OnLine = true
			playUrl.Width = width
			playUrl.Height = height
			playUrl.Ready = ready
			//update play_url set quality = 5 where content_type = 4 and width <= 1920;
			//update play_url set quality = 4 where content_type = 4 and width <= 1280;
			//update play_url set quality = 3 where content_type = 4 and width <= 720;
			//update play_url set quality = 2 where content_type = 4 and width <= 640;
			//update play_url set quality = 1 where content_type = 4 and width <= 480;
			if width <= 480{
				playUrl.Quality = 1
			}else if width <= 640{
				playUrl.Quality = 2
			}else if width <= 720{
				playUrl.Quality = 3
			}else if width <= 1280{
				playUrl.Quality = 4
			}else{
				playUrl.Quality = 5
			}
			if err = db.Save(&playUrl).Error ; err != nil{
				logger.Error(err)
				return
			}
		}else{
			playUrl.OnLine = false
			if err = db.Save(&playUrl).Error ; err != nil{
				logger.Error(err)
				return
			}
		}
	}
}


func GetStreamRation(url string)(uint32,uint32,int64,error){
	var name,videoCodec, audioCodec string
	var frameRate float32

	var bitrate, videoBitrate, audioBitrate, duration, width, height uint32
	var size uint64
	var err error

	t1 := time.Now()
	c1 := make(chan string, 1)
	go func() {
		name , frameRate, videoCodec, audioCodec, bitrate, videoBitrate, audioBitrate, duration, width, height, size, err = util.FfmpegVideoInfo(url)
		if err != nil{
			logger.Error(err)
			c1 <- "error"
			return
		}
		logger.Debug("name" ,":",name)
		logger.Debug("videoCodec" ,":",videoCodec)
		logger.Debug("audioCodec" ,":",audioCodec)
		logger.Debug("frameRate" ,":",frameRate)
		logger.Debug("bitrate" ,":",bitrate)
		logger.Debug("videoBitrate" ,":",videoBitrate)
		logger.Debug("audioBitrate" ,":",audioBitrate)
		logger.Debug("duration" ,":",duration)
		logger.Debug("width" ,":",width)
		logger.Debug("height" ,":",height)
		logger.Debug("size" ,":",size)

		if err == nil {
			c1 <- "success"
		}else{
			logger.Error(err)
			c1 <- "error"
		}
	}()
	select {
	case res := <-c1:
		if res == "success"{
			elapsed := time.Since(t1)
			logger.Debug("耗时:",elapsed)
			return width,height,int64(elapsed),nil
		}else{
			return 0,0,0,errors.New("获取视频信息失败!!!")
		}
	case <-time.After(time.Second * 15):
		return 0,0,0,errors.New("获取视频信息超过10秒，超时!!!")
	}

}