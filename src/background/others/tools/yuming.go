package main

import (
	"net/http"
	"background/common/logger"
	"background/others/config"
	"io/ioutil"
	"encoding/xml"
	"strings"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	GetStatus("ezhantao.com")
}
//
//type Prop struct {
//	ReturnCode        string    `xml:"returncode"`
//	key               string    `xml:"key"`
//	Original          string    `xml:"original"`
//}
//
//type Result struct {
//	Property Prop `xml:"property"`
//}


type DominData struct {
	XMLName           xml.Name `xml:"property"`
	Version     string   `xml:"version,attr"`
	Prop         Property `xml:"property"`
	Description string   `xml:",innerxml"`
}

type Property struct {
	XMLName    xml.Name `xml:"property"`
	ReturnCode        string    `xml:"returncode"`
	key               string    `xml:"key"`
	Original          string    `xml:"original"`
}

func GetStatus(url string){

	apiUrl := "http://panda.www.net.cn/cgi-bin/check.cgi?area_domain=" + url
	requ, err := http.NewRequest("GET", apiUrl, nil)

	resp, err := http.DefaultClient.Do(requ)
	if err != nil {
		logger.Debug("Proxy failed!")
		return
	}

	recv, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debug(string(recv))
	d := string(recv)
	d1 := strings.Replace(d,"gb2312","utf-8",-1)
	var result DominData
	err = xml.Unmarshal([]byte(d1), &result)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debug(result.Prop.Original)
}