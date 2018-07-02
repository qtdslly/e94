package main

import (
	"fmt"

	"background/newmovie/config"
	"background/common/logger"
	"background/common/util"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())
	for i := 1 ; i < 999; i++{
		//http://91.149.172.100:80/channel?id=94&u=demo&p=demo
		//http://61.58.60.230:9000/live/708.m3u8
		//url := "http://91.149.172.100:80/channel?id=" + fmt.Sprint(i) + "&u=demo&p=demo"
		//url := "http://61.58.60.230:9000/live/" + fmt.Sprint(i) + ".m3u8"
		//rtmp://wv4.tp33.net/sat/tv571
		//url := "rtmp://wv4.tp33.net/sat/tv" + fmt.Sprintf("%03d",i)
		//http://47.95.172.168/channel/2437.m3u8
		//url := "http://47.95.172.168/channel/" + fmt.Sprintf("%04d",i) + ".m3u8"
		//http://62.210.214.28/live?channelId=143&uid=171&deviceUser=seanmccarr&devicePass=120817
		//url := "http://62.210.214.28/live?channelId=" + fmt.Sprint(i) + "&uid=171&deviceUser=seanmccarr&devicePass=120817"
		//rtmp://edge1.everyon.tv:1935/etv1sb/phd563
		url := "http://223.110.243.142/PLTV/2510088/224/3221227" + fmt.Sprintf("%03d",i) + "/1.m3u8"
		if util.CheckStreamUrl(url){
			fmt.Println("SUCCESS" + url)
		}
	}
}


