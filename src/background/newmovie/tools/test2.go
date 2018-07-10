package main

import (
	"background/common/logger"
	"background/newmovie/config"
	"background/newmovie/service"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	word := "战狼1"
	_,title,score,area,description,actors,directors,thumb,pageUrl,publishDate := service.GetIqiyiVideoInfoByTitle(word)
	logger.Debug("豆瓣评分:",score)
	logger.Debug("截图:",thumb)
	logger.Debug("标题:",title)
	logger.Debug("导演:",directors)
	logger.Debug("主演:",actors)
	logger.Debug("描述:",description)
	logger.Debug("视频地址:",pageUrl)
	logger.Debug("地区:",area)
	logger.Debug("上映时间:",publishDate)

	logger.Debug("=========================================================================================")

	_,title,description,thumb,score,actors,directors,pageUrl = service.GetTencentVideoInfoByTitle(word)
	logger.Debug("豆瓣评分:",score)
	logger.Debug("截图:",thumb)
	logger.Debug("标题:",title)
	logger.Debug("导演:",directors)
	logger.Debug("主演:",actors)
	logger.Debug("描述:",description)
	logger.Debug("视频地址:",pageUrl)

	logger.Debug("=========================================================================================")

	_,title,description,actors,directors,thumb,pageUrl,publishDate = service.GetYoukuVideoInfoByTitle(word)
	logger.Debug("截图:",thumb)
	logger.Debug("标题:",title)
	logger.Debug("导演:",directors)
	logger.Debug("主演:",actors)
	logger.Debug("描述:",description)
	logger.Debug("视频地址:",pageUrl)
	logger.Debug("上映时间:",publishDate)
}
