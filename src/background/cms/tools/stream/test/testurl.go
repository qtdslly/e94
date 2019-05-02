package main

import (
	"fmt"

	"background/newmovie/config"
	"background/common/logger"
	"background/common/util"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())
	url := "http://ginocdn.bybzj.com:8092/wuma/20180101/hd_heyzo-1628/650kb/hls/index.m3u8"
	//url := "http://ginocdn.bybzj.com:8092/wuma/20180101/hd_heyzo-1628/650kb/hls/index.m3u8"
	if util.CheckStreamUrl(url,"aaa.jpg"){
		fmt.Println("SUCCESS" + url)
	}
}


