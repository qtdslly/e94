package script
/*http://www.meijutt.com*/

import (
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

	url1 := u.Host + u.Path
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
	redictUrl := GetRRmeijuValue("var redirecturl",data)
	mainStr := GetRRmeijuValue("var main",data)
	//mp4 := GetRRmeijuValue("mp4",data)
	realUrl := redictUrl + mainStr

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

	value = strings.Replace(value,"\r","",-1)
	value = strings.Replace(value,"\n","",-1)

	return value
}

