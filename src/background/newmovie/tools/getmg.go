package main

import (
	"background/newmovie/config"

	"github.com/PuerkitoBio/goquery"
	"background/common/logger"
	"strings"
	"fmt"
	"background/newmovie/model"
	"time"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"io/ioutil"
	"github.com/jinzhu/gorm"
	"flag"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())
	configPath := flag.String("conf", "../config/config.json", "Config file path")
	err := config.LoadConfig(*configPath)
	if err != nil {
		return
	}

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}

	i := 1
	for {
		url := "https://list.mgtv.com/3/a4-a3-------a7-1-" + fmt.Sprint(i) +  "--b1-.html?channelId=3"
		GetMgMovie(url,db)
		i = i + 1
		if i == 14{
			break
		}
	}
}


func GetMgMovie(url string,db *gorm.DB){
	query := GetMgPageInfo(url)

	if query != nil{
		FilterMgMovieInfo(query,db)
	}
}

type MgtvVideoInfo struct {
	Area	             string             `gorm:"area" json:"area"`
	Duration                 string        `gorm:"duration" json:"duration"`
	Release                 string        `gorm:"release" json:"release"`
}
type MgtvDataInfo struct {
	Info                 MgtvVideoInfo          `gorm:"info" json:"info"`
}
type MgtvData struct{
	Code                int             `gorm:"code" json:"code"`
	Data                 MgtvDataInfo    `gorm:"data" json:"data"`
	Msg                string             `gorm:"msg" json:"msg"`
	Seqid                string             `gorm:"seqid" json:"seqid"`
}

func FilterMgMovieInfo(document *goquery.Document,db *gorm.DB)(){
	movieDoc := document.Find("body").Find(".m-main").Find(".m-content").Find(".m-content-wrapper").Find(".m-result").Find(".m-result-list").Find("ul").Find("li")

	movieDoc.Each(func(i int, s *goquery.Selection) {
		title := s.Find(".u-title").Eq(0).Text()
		url,_ := s.Find(".u-title").Eq(0).Attr("href")
		if !strings.Contains(url,"http"){
			url = "http:" + url
		}
		thumb_y,_ := s.Find("img").Eq(0).Attr("src")
		if !strings.Contains(thumb_y,"http"){
			thumb_y = "http:" + thumb_y
		}

		score := s.Find("a").Eq(0).Find("em").Eq(0).Text()

		logger.Debug(title)
		logger.Debug(url)
		logger.Debug(thumb_y)
		logger.Debug(score)

		directors,actors,country,tags,description := GetMgInfoByUrl(url)

		logger.Debug(country)
		logger.Debug(score)
		logger.Debug(directors)
		logger.Debug(description)
		logger.Debug(tags)

		logger.Debug(url)
		tmp := strings.Replace(url,"http://www.mgtv.com/b/","",-1)
		tmp = strings.Replace(tmp,".html","",-1)
		videoIds := strings.Split(tmp,"/")
		if len(videoIds) != 2{
			logger.Error("获取videoid错误!!!")
			return
		}

		p := time.Now()
		t := fmt.Sprint(p.Unix()) + "000"
		apiurl := "https://pcweb.api.mgtv.com/movie/list?video_id=" + videoIds[1] + "&cxid=&version=5.5.35&callback=jQuery18205330103375947974_1528805597948&_support=10000000&_=" + t

		requ, err := http.NewRequest("GET",apiurl,nil)

		resp, err := http.DefaultClient.Do(requ)
		if err != nil {
			logger.Debug("Proxy failed!")
			return
		}

		recv,err := ioutil.ReadAll(resp.Body)
		if err != nil{
			logger.Error(err)
			return
		}

		data := string(recv)
		if data == ""{
			return
		}

		//jQuery18204076358277541021_1528805935312({"code":200,"data":{"info":{"area":"|内地,49|","duration":"20分钟","img":"https://3img.hitv.com/preview/internettv/sp_images/ott/2017/dianying/319142/20171024182204232-new.jpg_220x308.jpg","isvip":"0","kind":"|武侠,2831480|爱情,1275381|","lang":"普通话","playcnt":"954.1万","release":"2017-10-29","title":"忘忧镇","type":"0","url":"/b/319142/4159520.html","video_id":"4159520"},"list":[{"playcnt":"163.6万","title":"微电影《忘忧镇》互动版 陪赵丽颖林更新闯荡江湖","url":"/b/319142/4159520.html","video_id":"4159520"},{"playcnt":"147.7万","title":"微电影《忘忧镇》剧情版 赵丽颖林更新再续前缘","url":"/b/319142/4159632.html","video_id":"4159632"},{"playcnt":"23.3万","title":"微电影《忘忧镇》赵丽颖版 揭开神秘女侠的离奇身世","url":"/b/319142/4159880.html","video_id":"4159880"},{"playcnt":"394.5万","title":"微电影《忘忧镇》林更新版 林愈旧被追杀遇侠士相救","url":"/b/319142/4159960.html","video_id":"4159960"}],"series":[],"related":null,"short":[{"attr":"2","clip_id":"319142","img":"https://1img.hitv.com/preview/sp_images/2017/dianying/319142/4159519/20171103162504274.jpg_220x125.jpg","isnew":"0","playcnt":"36.6万","t1":"赵丽颖林更新命运归往何处？","t2":"00:55","t3":"微电影《忘忧镇》首发先导片 赵丽颖林更新命运归往何处？","url":"/b/319142/4159519.html","video_id":"4159519"},{"attr":"4","clip_id":"319142","img":"https://1img.hitv.com/preview/sp_images/2017/dianying/319142/4153244/20171030153935694.jpg_220x125.jpg","isnew":"0","playcnt":"35.2万","t1":"赵丽颖林更新出场配一脸","t2":"01:37","t3":"微电影《忘忧镇》首映礼 赵丽颖林更新出场配一脸","url":"/b/319142/4153244.html","video_id":"4153244"},{"attr":"4","clip_id":"319142","img":"https://1img.hitv.com/preview/sp_images/2017/dianying/319142/4153243/20171030153812978.jpg_220x125.jpg","isnew":"0","playcnt":"22.4万","t1":"赵丽颖林更新互动甜蜜飞","t2":"01:30","t3":"微电影《忘忧镇》首映礼 赵丽颖林更新互动甜蜜飞","url":"/b/319142/4153243.html","video_id":"4153243"},{"attr":"4","clip_id":"319142","img":"https://3img.hitv.com/preview/sp_images/2017/dianying/319142/4153242/20171030153701119.jpg_220x125.jpg","isnew":"0","playcnt":"18.5万","t1":"赵丽颖林更新热聊爆笑不止","t2":"00:45","t3":"微电影《忘忧镇》首映礼 赵丽颖林更新热聊爆笑不止","url":"/b/319142/4153242.html","video_id":"4153242"},{"attr":"2","clip_id":"319142","img":"https://3img.hitv.com/preview/sp_images/2017/dianying/319142/4153241/20171030153439856.jpg_220x125.jpg","isnew":"0","playcnt":"37.1万","t1":"白衣女侠仗义走江湖","t2":"03:58","t3":"微电影《忘忧镇》赵丽颖片花 白衣女侠仗义走江湖","url":"/b/319142/4153241.html","video_id":"4153241"},{"attr":"2","clip_id":"319142","img":"https://1img.hitv.com/preview/sp_images/2017/dianying/319142/4153240/20171030153132893.jpg_220x125.jpg","isnew":"0","playcnt":"17.5万","t1":"文弱书生的剑侠情缘","t2":"03:55","t3":"微电影《忘忧镇》林更新片花 文弱书生的剑侠情缘","url":"/b/319142/4153240.html","video_id":"4153240"},{"attr":"4","clip_id":"319142","img":"https://1img.hitv.com/preview/sp_images/2017/dianying/319142/4146469/20171024193130510.jpg_220x125.jpg","isnew":"0","playcnt":"20.9万","t1":"蒙面女侠赵丽颖的卖萌日记","t2":"00:58","t3":"微电影《忘忧镇》赵丽颖花絮 蒙面女侠的卖萌日记","url":"/b/319142/4146469.html","video_id":"4146469"},{"attr":"4","clip_id":"319142","img":"https://2img.hitv.com/preview/sp_images/2017/dianying/319142/4146468/20171024192411181.jpg_220x125.jpg","isnew":"0","playcnt":"20.4万","t1":"林大侠片场迷之日常","t2":"01:31","t3":"微电影《忘忧镇》林更新花絮 林大侠片场迷之日常","url":"/b/319142/4146468.html","video_id":"4146468"},{"attr":"3","clip_id":"319142","img":"https://0img.hitv.com/preview/sp_images/2017/dianying/319142/4146465/20171024193217572.jpg_220x125.jpg","isnew":"0","playcnt":"18.1万","t1":"林更新赵丽颖再续虐恋","t2":"00:16","t3":"微电影《忘忧镇》预告片 林更新赵丽颖再续虐恋","url":"/b/319142/4146465.html","video_id":"4146465"}]},"msg":"","seqid":"963835484c8a476eb6081db795136c75"})

		data = strings.Replace(data,"jQuery18205330103375947974_1528805597948(","",-1)
		data = data[0:len(data) - 1]

		//data, _ = service.DecodeToGBK(data)
		logger.Debug(data)

		var mgData MgtvData
		if err = json.Unmarshal([]byte(data), &mgData); err != nil {
			logger.Error(err)
			return
		}

		publishDate := mgData.Data.Info.Release
		duration := mgData.Data.Info.Duration
		duration = strings.Replace(duration,"分钟","",-1)
		logger.Debug(publishDate)
		logger.Debug(duration)

		var movie  model.Movie
		movie.Provider = "mgtv"
		movie.Actors = actors
		movie.Title = title
		movie.Description = description
		movie.Directors = directors
		movie.Url = url
		movie.PublishDate = publishDate
		movie.Score = score
		movie.ThumbY = thumb_y
		//movie.Year = year
		movie.Country = country
		movie.Duration = duration
		movie.Tags = tags
		now := time.Now()
		movie.CreatedAt = now
		movie.UpdatedAt = now
		if err := db.Where("title = ? and provider = ?",movie.Title,movie.Provider).First(&movie).Error ; err == gorm.ErrRecordNotFound{
			db.Create(&movie)
		}
	})
}


func GetMgInfoByUrl(url string)(string,string,string,string,string){
	query := GetMgPageInfo(url)

	base := query.Find(".play-container").Find(".play-primary").Find(".play-primary").Find(".v-panel").Find(".v-panel-box").Find(".v-panel-info").Find(".extend").Find(".v-panel-extend").Find(".v-panel-meta")

	base = base.Find("p")
	directors := base.Eq(0).Text()
	directors = strings.Replace(directors,"导演：","",-1)
	directors = strings.Replace(directors," ","",-1)
	directors = strings.Replace(directors,"\n","",-1)

	actors := base.Eq(1).Text()
	actors = strings.Replace(actors,"主演：","",-1)
	actors = strings.Replace(actors," ","",-1)
	actors = strings.Replace(actors,"\n","",-1)

	country := base.Eq(2).Text()
	country = strings.Replace(country,"地区：","",-1)
	country = strings.Replace(country," ","",-1)
	country = strings.Replace(country,"\n","",-1)


	tags := base.Eq(3).Text()
	tags = strings.Replace(tags,"类型：","",-1)
	tags = strings.Replace(tags," ","",-1)
	tags = strings.Replace(tags,"\n","",-1)

	description := base.Eq(4).Text()
	description = strings.Replace(description,"简介：","",-1)
	description = strings.Replace(description," ","",-1)
	description = strings.Replace(description,"\n","",-1)

	return directors,actors,country,tags,description
}

func GetMgPageInfo(url string)( *goquery.Document){
	logger.Debug(url)
	query, err := goquery.NewDocument(url)
	if err != nil {
		logger.Debug(url)
		logger.Error(err)
		return nil
	}
	return query
}