package main

import (
	"net/http"
	//"strings"
	"background/common/logger"
	"io/ioutil"
	"background/newmovie/config"

	"strings"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	apiUrl := "http://up.gitcms.com/webcloud/post.php"
	//postString := "{\"vid\":\"86217539-1-1-1-1-1-1\"}"

	requ, err := http.NewRequest("POST", apiUrl, strings.NewReader("vid=86217539-1-1-1-1-1-1"))
	//requ.Header.Add("Referer", "https://cloud.baidu.com/product/bcd/search.html?keyword=ezhantao")
	//requ.Header.Add("Host", "up.gitcms.com")
	//requ.Header.Add("Cookie", "time_http://fuli.zuida-youku-le.com/20180528/27035_ca5bf0af/index.m3u8=292.187499")
	//requ.Header.Add("Origin", "http://up.gitcms.com")
	requ.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//requ.Header.Add("Content-Type", "application/json, text/javascript, */*; q=0.01")
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
}
