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



//{"domainNames":[{"label":"ezhantao","tld":"com"}]}

type Domain struct {
	Label        string    `json:"label"`
	Tld	     string    `json:"tld"`
}
type DomainReq struct {
	DomainNames []Domain  `json:domainNames`
}


func GetStatus(url string){

	apiUrl := "https://cloud.baidu.com/api/bcd/search/status" + url
	requ, err := http.NewRequest("POST", apiUrl, nil)
	requ.Header.Add("Host", "cloud.baidu.com")
	requ.Header.Add("Referer", "https://cloud.baidu.com/product/bcd/search.html?keyword=ezhantao")
	requ.Header.Add("Host", "cloud.baidu.com")
	requ.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3278.0 Safari/537.36")

	req, err := http.NewRequest("POST", apiUrl, strings.NewReader("name=cjb"))
	if err != nil {
		// handle error
	}

	var dom Domain
	dom.Label = "ezhantao"
	dom.Tld = "com"
	domr := DomainReq{dom}
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