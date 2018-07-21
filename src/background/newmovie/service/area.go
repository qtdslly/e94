package service

import (
	"os"
	"bufio"
	"io"
	"strings"
	"background/common/logger"

)
var AREA []string

func SetArea(file string){
	f, err := os.Open(file)
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

		fields := strings.Split(line, "|")
		area := fields[1]
		area = strings.Replace(area,"区","",-1)
		area = strings.Replace(area,"县","",-1)
		area = strings.Replace(area,"自治","",-1)
		area = strings.Replace(area,"市","",-1)
		area = strings.Replace(area,"前旗","",-1)
		area = strings.Replace(area,"中旗","",-1)
		area = strings.Replace(area,"后旗","",-1)
		area = strings.Replace(area,"旗","",-1)
		area = strings.Replace(area,"州","",-1)

		AREA = append(AREA,area)
	}
}
