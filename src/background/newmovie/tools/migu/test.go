package main

import (
	"fmt"
	"background/newmovie/service/script"
	"background/newmovie/config"
	"background/common/logger"
	"flag"
	"background/common/constant"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())
	configPath := flag.String("conf", "../config/config.json", "Config file path")
	err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Error(err)
		return
	}

	url := "http://www.miguvideo.com/wap/resource/pc/detail/miguplay.jsp?cid=608653476"
	logger.Debug(url)
	realUrl := script.GetMiguRealPlayUrl(constant.MediaTypeEpisode,url)
	fmt.Println(realUrl)
}