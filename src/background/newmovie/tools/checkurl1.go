package main

import (
	"fmt"

	"background/newmovie/config"
	"background/common/logger"
	"background/common/util"
	"os"
	"bufio"
	"io"
	"strings"
)
func main(){
	logger.SetLevel(config.GetLoggerLevel())
	var err error
	f, err := os.Open("/root/Git/e94/src/background/newmovie/tools/stream.txt")
	if err != nil {
		logger.Error(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		line = strings.Replace(line, "\n", "", -1)

		logger.Debug(line)

		var title string
		var url string
		fields := strings.Split(line, "|")
		if len(fields) == 2{
			title = fields[0]
			url = fields[1]
		}else{
			url = fields[0]
		}
		if util.CheckStreamUrl(url){
			fmt.Println("SUCCESS" + title + "\t" + url)
		}
	}
}

