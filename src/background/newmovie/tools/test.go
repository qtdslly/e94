package main

import (
	"background/newmovie/service"
	"background/newmovie/config"

	"background/common/logger"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	url := service.GetRealUrl("youku","http://v.youku.com/v_show/id_XMzU0ODk0MzQ0MA==.html",service.GetJsCode())
	logger.Debug(url)
}
