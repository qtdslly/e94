package service

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"background/common/logger"
	"1/github.com/pkg/errors"
)

func GetIqiyiVideoInfoByTitle(word string)(error,string,string,string,string,string,string,string,string,string){
	var title,score,area,description,actors,directors,thumb,pageUrl,publishDate string
	/*
	http://so.iqiyi.com/so/q_%E6%88%98%E7%8B%BC2?source=input
	*/

	parWord := url.QueryEscape(word)

	apiUrl := "http://so.iqiyi.com/so/q_" + parWord + "?source=input"

	//logger.Debug(apiUrl)

	query, err := goquery.NewDocument(apiUrl)
	if err != nil {
		logger.Debug(apiUrl)
		logger.Error(err)
		return err,title,score,area,description,actors,directors,thumb,pageUrl,publishDate
	}


	//logger.Debug(query.Html())

	base := query.Find(".mod_result").Eq(0).Find(".mod_result_list").Eq(0).Find(".list_item").Eq(0)

	//logger.Debug(base.Html())

	title = query.Find(".result_info").Eq(0).Find(".result_title").Eq(0).Find("a").Eq(0).Text()

	thumb,exists := base.Find(".figure").Eq(0).Find("img").Eq(0).Attr("src")
	if !exists{
		thumb = ""
	}else{
		if !strings.Contains(thumb,"http"){
			thumb = "http:" + thumb
		}
	}

	pageUrl,exists = base.Find(".figure").Eq(0).Attr("href")
	if !exists{
		return errors.New("网页地址获取失败"),title,score,area,description,actors,directors,thumb,pageUrl,publishDate
	}

	score = base.Find(".result_info").Eq(0).Find(".result-info-score").Text()

	others := base.Find(".result_info").Eq(0).Find(".info_item").Find(".result_info_cont")

	others.Each(func(i int, s *goquery.Selection) {

		if strings.Contains(s.Text(),"上映时间"){
			publishDate = strings.Replace(s.Text(),"上映时间:","",-1)
		}else if strings.Contains(s.Text(),"主演"){
			actors = strings.Replace(s.Text(),"主演:","",-1)
		}else if strings.Contains(s.Text(),"导演"){
			directors = strings.Replace(s.Text(),"导演:","",-1)
		}else if strings.Contains(s.Text(),"地区"){
			area = strings.Replace(s.Text(),"地区:","",-1)
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
	publishDate = strings.Replace(publishDate,"\n","",-1)
	area = strings.Replace(area,"\n","",-1)
	score = strings.Replace(score,"\n","",-1)

	thumb = strings.Replace(thumb," ","",-1)
	title = strings.Replace(title," ","",-1)
	directors = strings.Trim(directors," ")
	actors = strings.Trim(actors," ")
	description = strings.Replace(description," ","",-1)
	pageUrl = strings.Replace(pageUrl," ","",-1)
	publishDate = strings.Replace(publishDate," ","",-1)
	area = strings.Replace(area," ","",-1)
	score = strings.Replace(score," ","",-1)

	thumb = strings.Replace(thumb,"	","",-1)
	title = strings.Replace(title,"	","",-1)
	directors = strings.Trim(directors,"	")
	actors = strings.Trim(actors,"	")
	description = strings.Replace(description,"	","",-1)
	pageUrl = strings.Replace(pageUrl,"	","",-1)
	publishDate = strings.Replace(publishDate,"	","",-1)
	area = strings.Replace(area,"	","",-1)
	score = strings.Replace(score," ","",-1)

	actors = TrimSpace(actors)
	directors = TrimSpace(directors)

	actors = strings.Replace(actors," ","/",-1)
	directors = strings.Replace(directors," ","/",-1)
	description = description[0:strings.LastIndex(description,"详细")]
	return nil,title,score,area,description,actors,directors,thumb,pageUrl,publishDate
}

func TrimSpace(word string)string{
	if !strings.Contains(word,"  "){
		return word
	}
	word = strings.Replace(word,"  "," ",-1)
	return TrimSpace(word)
}