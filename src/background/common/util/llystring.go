package util

import (
	"strings"
	"time"
	//"fmt"
	"os/exec"
	"background/common/logger"
)



//保留数字和字母以及+-符号
func TrimChinese(title string)(string){
	if !strings.Contains(title,"CCTV"){
		return title
	}
	if strings.Contains(title,"证券资讯") || strings.Contains(title,"世界地理") ||strings.Contains(title,"书画") ||strings.Contains(title,"娱乐") ||strings.Contains(title,"怀旧剧场") || strings.Contains(title,"电视指南") ||strings.Contains(title,"第一剧场") || strings.Contains(title,"风云剧场") ||strings.Contains(title,"风云足球") ||strings.Contains(title,"风云音乐") ||strings.Contains(title,"高尔夫") {
		title = strings.Trim(title,"-")
		return title
	}
	//48-57 45 43 65-90 97-122
	rTitle := ([]rune)(title)
	result := ""
	for _, m := range rTitle {
		if m == 43 || m == 45 || (m >= 48 && m <= 57) || (m >= 65 && m <=90) || (m >= 97 && m <= 122){
			result += string(m)
		}
	}
	return result
}


func CheckStream(url string,jpgName string)bool{
	for i := 0 ; i < 3 ; i++{
		if CheckStreamUrl(url,jpgName){
			return true
		}
	}
	return false
}


func CheckStreamUrl(url string,jpgName string)bool{
	c2 := make(chan string, 1)
	//var ffmpegAddr string = "ffmpeg"
	go func() {
		//logger.Debug(ffmpegAddr)
		logger.Debug(url)

		var cmdStr string = "C:/ff/bin/ffmpeg -i " + url + " -y -s 320x240 -vframes 1 " + jpgName
		//cmdStr := fmt.Sprintf("ffmpeg -i %s -y -s 320x240 -vframes 1 %s", url,jpgName)
		logger.Debug(cmdStr)
		time.Sleep(time.Second * 5)
		cmd := exec.Command("bash", "-c", "C:/ff/bin/ffmpeg -i " + url + " -y -s 320x240 -vframes 1 " + jpgName)

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
			return true
		}else{
			return false
		}
	case <-time.After(time.Second * 30):
		return false
	}

	return false
}