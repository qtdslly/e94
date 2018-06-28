package main

import (
	"background/common/logger"
	"background/newmovie/service/script"
	"background/newmovie/config"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	realUrl := script.GetMiguRealPlayUrl(2,"http://www.miguvideo.com/wap/resource/pc/detail/miguplay.jsp?cid=617379229")
	logger.Debug(realUrl)
}