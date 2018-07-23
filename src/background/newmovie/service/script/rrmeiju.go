package script
/*http://www.meijutt.com*/


import (
	"fmt"
	"background/common/logger"
	"io/ioutil"
	"net/http"
	"strings"
	"github.com/tidwall/gjson"
	"github.com/PuerkitoBio/goquery"
	"net/url"
)
func GetRRmeijuRealPlayUrl(apiurl string)(string){

	query, err := goquery.NewDocument(apiurl)
	if err != nil {
		logger.Debug(apiurl)
		logger.Error(err)
		return ""
	}

	base := query.Find("script")
	html := ""
	base.Each(func(i int, s *goquery.Selection) {

		if strings.Contains(s.Text(),"var player_data=") {
			html = s.Text()
		}

	})
	if html == ""{
		return ""
	}

	//logger.Debug(html)
	html = strings.Replace(html,"var player_data=","",-1)
	//logger.Debug(html)

	apiUrl1 := gjson.Get(html, "url")

	if !apiUrl1.Exists() {
		return ""
	}


	u, err := url.Parse(apiUrl1.String())

	url1 := fmt.Sprintf(`%s%s`, u.Host, u.Path)

	//logger.Debug(url1)
	requ, err := http.NewRequest("GET", url1,nil)

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

	data := string(recv)
	logger.Debug(data)
	redictUrl := GetRRmeijuValue("redirecturl",data)
	mainStr := GetRRmeijuValue("main",data)
	//redictUrl := "http://vs1.baduziyuan.com"
	//mainStr := "/20180412/g2b3HvWx/index.m3u8?sign=3dd8344fb2d3567a4d400213938a6d9269450666637dca14ea1599d8cb583fe967e77b29236ea6465d4ca1c5e723f0caf42f8b7d51aa92f5180ee0fa103b88a4"
	//mp4 := GetRRmeijuValue("mp4",data)
	logger.Debug(redictUrl)
	logger.Debug(mainStr)
	//logger.Debug(redictUrl)
	//logger.Debug(mp4)
	//realUrl := redictUrl + mainStr
	realUrl := fmt.Sprintf("%s%s", redictUrl, mainStr)
	logger.Debug(realUrl)
	return realUrl
}

func GetRRmeijuValue(label ,pageinfo string)(string){
	ss := strings.Split(pageinfo,"\n")
	var value string
	for _ , s:= range ss{
		if strings.Contains(s,label){
			start := strings.Index(s,label)
			if start >= 0{
				value = s[start + len(label) + 3:]
				value = strings.Replace(value,"\"","",-1)
				value = strings.Replace(value,";","",-1)
				break
			}

		}
	}
	return value
}

