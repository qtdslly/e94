package script
/*www.hanju.cc/hanju/*/


import (
	"background/common/logger"
	"io/ioutil"
	"net/http"
	"strings"
)
func GetHanjuRealPlayUrl(url string)(string){

	requ, err := http.NewRequest("GET", url,nil)

	resp, err := http.DefaultClient.Do(requ)
	if err != nil {
		logger.Error(err)
		return ""
	}

	recv,err := ioutil.ReadAll(resp.Body)
	if err != nil{
		logger.Error(err)
		return ""
	}
	//var vid='https://www2.yuboyun.com/hls/2018/07/08/nVMgCubW/playlist.m3u8';

	data := string(recv)
	start := strings.Index(data,"var vid='")
	end := strings.Index(data,".m3u8");
	if start == -1 || end == -1{
		return ""
	}

	return data[start+9:end+5]
}

