package util

import (
	"strings"
	"time"
	"fmt"
	"os/exec"
	"background/common/logger"
)



//保留数字和字母以及+-符号
func TrimChinese(title string)(string){
	if !strings.Contains(title,"CCTV"){
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



func CheckStreamUrl(url string)bool{
	c2 := make(chan string, 1)
	ffmpegAddr := "/usr/bin/ffmpeg"
	go func() {
		cmdStr := fmt.Sprintf("%s -i '%s' -y -s 320x240 -vframes 1 aaa.jpg", ffmpegAddr, url)
		fmt.Println(cmdStr)
		cmd := exec.Command("bash", "-c", cmdStr)

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
	case <-time.After(time.Second * 10):
		return false
	}

	return false
}