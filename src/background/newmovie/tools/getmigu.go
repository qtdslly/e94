package main

import (
	"background/newmovie/config"
	"background/common/util"
	"background/common/logger"
	"strings"
	"fmt"
	"github.com/tidwall/gjson"

	"background/newmovie/model"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"io/ioutil"
	"github.com/jinzhu/gorm"
	"flag"
	"background/common/constant"
	"strconv"
)

func main(){
	logger.SetLevel(0)
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
	db.LogMode(true)

	model.InitModel(db)
	i := 1
	for {
		url := "http://www.miguvideo.com/wap/resource/pc/data/filmLibraryData.jsp?type=1&searchType=1000&searchContentType=&searchYear=&searchArea=&searchLimit=1002601&pageSize=25&currentPage=" + fmt.Sprint(i) + "&order=2&searchShape=%E5%85%A8%E7%89%87"
		GetMiguMovie(url,db)
		i = i + 1
		if i == 632{
			break
		}
	}
}


func GetMiguMovie(url string,db *gorm.DB){
	requ, err := http.NewRequest("GET", url,nil)
	requ.Header.Add("Host", "www.miguvideo.com")
	requ.Header.Add("Referer", "http://www.miguvideo.com/wap/resource/pc/list/filmLibrary.jsp?typeName=%E7%94%B5%E5%BD%B1")
	requ.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
	requ.Header.Add("Cookie", "TouristAccess=2018062510233616363F42A92DCABD047EF0D05BFD2D12; BIGipServerFarm-WAP-PBS-WIFI=2199242762.20480.0000; WT_FPC=id=21ed00571d0ed8267c71529894464452:lv=1530069203845:ss=1530069141539; JSESSIONID=2B148860CB7F4B86691323D0209ECC74")
	requ.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")

	resp, err := http.DefaultClient.Do(requ)
	if err != nil {
		logger.Debug("err")
		return
	}

	recv,err := ioutil.ReadAll(resp.Body)
	if err != nil{
		logger.Error(err)
		return
	}
	logger.Debug(url)
	data := string(recv)
	logger.Debug(data)
	content := gjson.Get(data, "content")

	if content.Exists() {
		videos := content.Array()
		for _, v := range videos {
			tx := db.Begin()
			var video model.Video
			video.Title = v.Get("title").String()
			num,_ := strconv.Atoi(v.Get("nf").String())
			video.Year = uint32(num)
			video.Description = v.Get("abstruct").String()
			video.ThumbY = v.Get("imgsrcV").String()
			video.ThumbX = v.Get("imgsrcH").String()
			video.Category = v.Get("DisplayName").String()

			score,_ := strconv.ParseFloat(v.Get("score").String(),10)

			video.Score = score
			video.Directors = v.Get("dy").String()
			video.Tags = v.Get("lx").String()
			video.Actors = v.Get("zy").String()
			video.Country = v.Get("dq").String()
			video.CurrentEpisode = uint32(1)
			video.TotalEpisode = uint32(1)
			video.OnLine = constant.OnlineTypeTrue
			video.Pinyin = util.TitleToPinyin(video.Title)
			video.Sort = 0
			if strings.Contains(video.Country,"内地") || strings.Contains(video.Country,"大陆") || strings.Contains(video.Country,"中国") || strings.Contains(video.Country,"香港") || strings.Contains(video.Country,"台湾"){
				video.Language = "汉语"
			}else if(strings.Contains(video.Country,"法国")){
				video.Language = "法语"
			}else if(strings.Contains(video.Country,"俄罗斯")){
				video.Language = "俄语"
			}else if(strings.Contains(video.Country,"日本")){
				video.Language = "日语"
			}else if(strings.Contains(video.Country,"韩国")){
				video.Language = "韩语"
			}else if(strings.Contains(video.Country,"意大利")){
				video.Language = "意大利语"
			}else if(strings.Contains(video.Country,"印度")){
				video.Language = "印地语"
			}else{
				video.Language = "英语"
			}

			if err = db.Where("title = ? and year = ?",video.Title,video.Year).First(&video).Error ; err == gorm.ErrRecordNotFound{
				if err = db.Create(&video).Error ; err != nil{
					logger.Error(err)
					tx.Rollback()
					return
				}
			}else{
				if err = db.Save(&video).Error ; err != nil{
					logger.Error(err)
					tx.Rollback()
					return
				}
			}

			var episode model.Episode
			episode.Title = video.Title
			episode.Pinyin = video.Pinyin
			episode.Description = video.Description
			num,_ = strconv.Atoi(v.Get("timelength").String())
			episode.Duration = uint32(num)
			episode.Number = "1"
			episode.Sort = uint32(1)
			episode.Score = video.Score
			episode.ThumbOttX = video.ThumbX
			episode.ThumbY = video.ThumbY
			episode.VideoId = video.Id
			if err = db.Where("video_id = ?",video.Id).First(&episode).Error ; err == gorm.ErrRecordNotFound{
				if err = db.Create(&episode).Error ; err != nil{
					logger.Error(err)
					tx.Rollback()
					return
				}
			}else{
				if err = db.Save(&episode).Error ; err != nil{
					logger.Error(err)
					tx.Rollback()
					return
				}
			}

			var playUrl model.PlayUrl
			playUrl.Sort = 1
			playUrl.ContentId = episode.Id
			playUrl.ContentType = constant.MediaTypeEpisode
			playUrl.OnLine = constant.OnlineTypeTrue
			playUrl.Provider = constant.ContentProviderMigu
			playUrl.Quality = 3
			playUrl.Title = video.Title
			playUrl.Url = "http://www.miguvideo.com/wap/resource/pc/detail/miguplay.jsp?cid=" + v.Get("cententId").String()
			if err = db.Where("content_type = ? and content_id = ? and provider = ?",playUrl.ContentType,playUrl.ContentId,playUrl.Provider).First(&playUrl).Error ; err == gorm.ErrRecordNotFound{
				if err = db.Create(&playUrl).Error ; err != nil{
					logger.Error(err)
					tx.Rollback()
					return
				}
			}else{
				if err = db.Save(&playUrl).Error ; err != nil{
					logger.Error(err)
					tx.Rollback()
					return
				}
			}
			tx.Commit()
		}
	}
}