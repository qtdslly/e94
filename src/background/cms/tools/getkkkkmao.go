package main

import (
	"background/common/logger"
	"background/newmovie/config"
	//"github.com/jinzhu/gorm"
	"github.com/PuerkitoBio/goquery"
)

var BASE_URL string = "http://www.kkkkmao.com"

func main(){
	logger.SetLevel(config.GetLoggerLevel())
	GetKkkkMaoInfo("科幻")
}

func GetKkkkMaoInfo(movieType string){
	query := GetPageInfo(movieType)

	if query != nil{
		FilterMovieInfo(query)
	}
}

func FilterMovieInfo(document *goquery.Document){
	movieDoc := document.Find("#letter-focus").Find(".letter-focus-item").Find("dd")

	movieDoc.Each(func(i int, s *goquery.Selection) {
		a := s.Find("a").Eq(0)
		name := a.Text()
		url,_ := a.Attr("href")
		url = BASE_URL + url
		logger.Debug(name,":",url)
	})


}


func GetPageInfo(movieType string)( *goquery.Document){
	var apiurl string
	if movieType == "科幻"{
		apiurl = "http://www.kkkkmao.com/Sciencefiction/"
	}else if movieType == "战争"{
		apiurl = "http://www.kkkkmao.com/War/"
	}
	query, err := goquery.NewDocument(apiurl)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return query
}
