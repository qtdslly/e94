package main

import (
	"fmt"
	"background/newmovie/config"

	"background/common/logger"
	"os/exec"
	//"os"
	//"bufio"
	//"io"
	//"strings"
)
var ffmpegAddr string = "f:/ffmpeg/bin/ffmpeg"
func main(){
	logger.SetLevel(config.GetLoggerLevel())
	//var err error
	//f, err := os.Open("f:/Git/e94/src/background/newmovie/tools/stream.txt")
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()
	//
	//rd := bufio.NewReader(f)
	//for {
	//	line, err := rd.ReadString('\n')
	//	if err != nil || io.EOF == err {
	//		break
	//	}
	//	line = strings.Replace(line, "\r\n", "", -1)
	//	logger.Debug(line)
	//
	//	fields := strings.Split(line,",")
	//	logger.Print(fields)
	//	title := fields[0]
	//	url := fields[1]
	//	if CheckUrl(url){
	//		fmt.Println("USE",title,"\t",url)
	//	}
	//}

	CheckUrl("http://live01.hebtv.com/channels/hebtv/video_channel_03/m3u8:800k")
}

func CheckUrl(url string)bool{
	cmdStr := fmt.Sprintf("%s -i \"%s\" -y -s 320x240 -vframes 1 \"%s.jpg\"", ffmpegAddr, url, "f:/Git/e94/src/background/newmovie/tools/a")
	logger.Debug(cmdStr)
	cmd := exec.Command("bash", "-c", cmdStr)

	if err := cmd.Run(); err == nil {
		logger.Debug("success")
		return true
	}else{
		logger.Error(err)
	}
	return false
}

