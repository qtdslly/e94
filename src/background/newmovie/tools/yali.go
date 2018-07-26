package main

import (
	"net/http"
	"background/common/logger"
	"io/ioutil"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"background/newmovie/config"
	"sync"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	var process *sync.Mutex
	process = new(sync.Mutex)
	var Count int = 0

	for{
		for{
			if Count > 10{
				time.Sleep(time.Millisecond * 100)
			}else{
				break
			}
		}

		go func(){
			process.Lock()
			Count++
			process.Unlock()
			apiUrl := "http://www.ezhantao.com:16882/cms/page"
			requ, err := http.NewRequest("GET", apiUrl, nil)
			requ.Header.Add("app_version", "1.0.0")
			//requ.Header.Add("Referer", "https://cloud.baidu.com/product/bcd/search.html?keyword=ezhantao")
			requ.Header.Add("app_key", "5Hwf3G51ZUPT")
			requ.Header.Add("installation_id","1807150014158383")
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

			process.Lock()
			Count--
			process.Unlock()
		}()
	}

	for{
		time.Sleep(time.Minute)
	}
}