package service

import (
	"net/url"
	"time"
	"fmt"
	 "math/rand"

	"github.com/PuerkitoBio/goquery"
	"background/common/logger"
	"github.com/tidwall/gjson"

	"strings"
	"io"
	"github.com/pkg/errors"
)

func GetYoukuVideoInfoByTitle(word string)(error,string,string,string,string,string,string,string){
	var title,description,actors,directors,thumb,pageUrl,publishDate string
	/*
	http://so.youku.com/search_video/q_%E7%BA%A2%E6%B5%B7%E8%A1%8C%E5%8A%A8?spm=a2h0k.11417342.searcharea.dbutton&_t=1531121505235
	http://so.youku.com/search_video/q_红海行动?spm=a2h0k.11417342.searcharea.dbutton&_t=1531121505235
	*/

	parWord := url.QueryEscape(word)

	now := time.Now().Unix()

	apiUrl := "http://so.youku.com/search_video/q_" + parWord + "?spm=a2h0k.11417342.searcharea.dbutton&_t=" + fmt.Sprint(now) + fmt.Sprintf("%03d",rand.Intn(1000))

	//logger.Debug(apiUrl)

	query, err := goquery.NewDocument(apiUrl)
	if err != nil {
		logger.Debug(apiUrl)
		logger.Error(err)
		return err,title,description,actors,directors,thumb,pageUrl,publishDate
	}

	//logger.Debug(query.Html())

	//title := query.Find("title").Eq(0).Text()
	base := query.Find("script")
	html := ""
	base.Each(func(i int, s *goquery.Selection) {

		if strings.Contains(s.Text(),"导演") && strings.Contains(s.Text(),"主演") && strings.Contains(s.Text(),"上映时间") {
			html = s.Text()
		}

	})
	if html == ""{
		return errors.New("接口调用错误"),title,description,actors,directors,thumb,pageUrl,publishDate
	}

	html = strings.Replace(html,"bigview.view(","",-1)
	html = strings.Replace(html,")","",-1)
	//logger.Debug(html)

	accurate := gjson.Get(html, "html")

	if !accurate.Exists() {
		return errors.New("接口调用错误"),title,description,actors,directors,thumb,pageUrl,publishDate
	}

	values := accurate.String()

	var r io.Reader = strings.NewReader(values)

	query,err = goquery.NewDocumentFromReader(r)
	if err != nil{
		logger.Error(err)
		return err,title,description,actors,directors,thumb,pageUrl,publishDate
	}

	doc := query.Find(".sk-result-list").Eq(0).Find(".sk-mod").Eq(0)
	//logger.Debug(doc.Html())

	thumb,exists := doc.Find(".mod-left").Eq(0).Find(".pack-cover").Eq(0).Find("img").Eq(0).Attr("src")
	if !exists{
		thumb = ""
	}else{
		if !strings.Contains(thumb,"http"){
			thumb = "http:" + thumb
		}
	}

	title = doc.Find(".mod-main").Eq(0).Find(".mod-header").Eq(0).Find(".spc-lv-1").Eq(0).Text()

	pageUrl,exists = doc.Find(".mod-left").Eq(0).Find(".sk-pack").Eq(0).Attr("href")
	if !exists{
		return errors.New("网页地址获取失败"),title,description,actors,directors,thumb,pageUrl,publishDate
	}

	others := doc.Find(".mod-main").Eq(0).Find(".mod-info").Eq(0).Find("span")

	others.Each(func(i int, s *goquery.Selection) {

		if strings.Contains(s.Text(),"上映时间"){
			publishDate = strings.Replace(s.Text(),"上映时间:","",-1)
		}else if strings.Contains(s.Text(),"主演"){
			actors = strings.Replace(s.Text(),"主演:","",-1)
		}else if strings.Contains(s.Text(),"导演"){
			directors = strings.Replace(s.Text(),"导演:","",-1)
		}else if strings.Contains(s.Text(),"简介"){
			description = strings.Replace(s.Text(),"简介:","",-1)
		}
	})

	thumb = strings.Replace(thumb,"\n","",-1)
	title = strings.Replace(title,"\n","",-1)
	directors = strings.Replace(directors,"\n","",-1)
	actors = strings.Replace(actors,"\n","",-1)
	description = strings.Replace(description,"\n","",-1)
	pageUrl = strings.Replace(pageUrl,"\n","",-1)
	publishDate = strings.Replace(publishDate,"	","",-1)

	thumb = strings.Replace(thumb," ","",-1)
	title = strings.Replace(title," ","",-1)
	directors = strings.Replace(directors," ","",-1)
	actors = strings.Replace(actors," ","",-1)
	description = strings.Replace(description," ","",-1)
	pageUrl = strings.Replace(pageUrl," ","",-1)
	publishDate = strings.Replace(publishDate,"	","",-1)

	thumb = strings.Replace(thumb,"	","",-1)
	title = strings.Replace(title,"	","",-1)
	directors = strings.Replace(directors,"	","",-1)
	actors = strings.Replace(actors,"	","",-1)
	description = strings.Replace(description,"	","",-1)
	pageUrl = strings.Replace(pageUrl,"	","",-1)
	publishDate = strings.Replace(publishDate,"	","",-1)

	return nil,title,description,actors,directors,thumb,pageUrl,publishDate
}