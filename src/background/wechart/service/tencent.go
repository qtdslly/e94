package service

import (
	"github.com/PuerkitoBio/goquery"
	"background/common/logger"
	"strings"
	"net/url"
	"github.com/pkg/errors"
)

func GetTencentVideoInfoByTitle(word string)(error,string,string,string,string,string,string,string){

	var title,description,thumb,score,actors,directors,pageUrl string
	parWord := url.QueryEscape(word)

	apiUrl := "https://v.qq.com/x/search/?q=" + parWord + "&stag=0&ses=qid%3DG8EYYycvC2-IjF4f_Gto4l_eUJOY-O5yO64hXeQwWcLSjL6j7kIRxQ"

	query, err := goquery.NewDocument(apiUrl)
	if err != nil {
		logger.Debug(apiUrl)
		logger.Error(err)
		return err,title,description,thumb,score,actors,directors,pageUrl
	}

	base := query.Find(".result_item").Eq(0)

	//logger.Debug(base.Html())

	score = base.Find(".result_score").Eq(0).Text()
	thumb,exists := base.Find("img").Eq(0).Attr("src")
	if !exists{
		thumb = ""
	}else{
		if !strings.Contains(thumb,"http"){
			thumb = "http:" + thumb
		}
	}
	title = base.Find(".result_title").Eq(0).Find("a").Eq(0).Text()
	language := base.Find(".result_title").Eq(0).Find("a").Eq(0).Find("span").Eq(0).Text()
	videoType := base.Find(".result_title").Eq(0).Find("a").Eq(0).Find("span").Eq(1).Text()
	title = strings.Replace(title,language,"",-1)
	title = strings.Replace(title,videoType,"",-1)

	dirs := base.Find("a[_stat='video:poster_v_导演']")
	dirs.Each(func(i int, s *goquery.Selection) {

		directors += s.Text() + "|"

	})

	if len(directors) > 0{
		directors = directors[0:len(directors) - 1]
	}

	acts := base.Find("a[_stat='video:poster_v_主演']")
	acts.Each(func(i int, s *goquery.Selection) {

		actors += s.Text() + "/"

	})

	if len(actors) > 0{
		actors = actors[0:len(actors) - 1]
	}

	pageUrl,exists = base.Find(".result_title").Eq(0).Find("a").Eq(0).Attr("href")
	if !exists{
		return errors.New("网页地址获取失败"),title,description,thumb,score,actors,directors,pageUrl
	}

	description = base.Find(".desc_text").Eq(0).Text()
	description = description[0:strings.LastIndex(description,"详细")]
	description = strings.Replace(description,"\n","",-1)

	return nil,title,description,thumb,score,actors,directors,pageUrl
}
