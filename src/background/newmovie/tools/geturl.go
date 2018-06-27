package main

import (
	"background/common/logger"
	"background/newmovie/service/script"
	"background/newmovie/config"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	realUrl := script.GetIqiyiRealPlayUrl("http://www.iqiyi.com/v_19rrhc74rc.html")
	logger.Debug(realUrl)
}