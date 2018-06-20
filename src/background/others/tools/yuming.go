package main

import (
	"net/http"
	"background/common/logger"
	"background/others/config"
	"io/ioutil"
	"strings"
	"net/url"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	GetStatus("ezhantao.com")
}



//{"domainNames":[{"label":"ezhantao","tld":"com"}]}

//type Domain struct {
//	Label        string    `json:"label"`
//	Tld	     string    `json:"tld"`
//}
//type DomainReq struct {
//	DomainNames []Domain  `json:domainNames`
//}


func GetStatus(url string){

	apiUrl := "https://cloud.baidu.com/api/bcd/search/status" + url
	postValue := url.Values{
		"label": {"ezhantao"},
		"tld": {"com"},
	}
	postString := postValue.Encode()
	requ, err := http.NewRequest("POST", apiUrl, strings.NewReader(postString))
	requ.Header.Add("Host", "cloud.baidu.com")
	requ.Header.Add("Referer", "https://cloud.baidu.com/product/bcd/search.html?keyword=ezhantao")
	requ.Header.Add("Host", "cloud.baidu.com")
	requ.Header.Add("Content-Type", "application/json;charset=UTF-8")
	requ.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3278.0 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(requ)
	if err != nil {
		logger.Error(err)
		return
	}

	defer resp.Body.Close()

	recv, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debug(string(recv))


	//whois https://cloud.baidu.com/api/bcd/whois/detail

}