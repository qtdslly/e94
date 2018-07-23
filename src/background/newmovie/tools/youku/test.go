package main

import (
	"fmt"
	"background/newmovie/service/script"
	"background/newmovie/config"
	"background/common/logger"
	"flag"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())
	configPath := flag.String("conf", "../config/config.json", "Config file path")
	err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Error(err)
		return
	}

	url := "http://www.iqiyi.com/v_19rrj6udbc.html#vfrm=2-4-0-1"
	logger.Debug(url)
	realUrl := script.GetIqiyiRealPlayUrl(url)
	fmt.Println(realUrl)
}