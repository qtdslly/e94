package main

import (
	"background/shortvideo/config"
	"background/common/logger"
	"background/common/constant"

	"background/shortvideo/model"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"flag"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
	"strconv"
)

func main() {
	logger.SetLevel(config.GetLoggerLevel())

	configPath := flag.String("conf", "../config/config.json", "Config file path")

	flag.Parse()

	logger.SetLevel(config.GetLoggerLevel())

	err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Error("Config Failed!!!!", err)
		return
	}

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}

	db.LogMode(true)
	model.InitModel(db)
	logger.SetLevel(config.GetLoggerLevel())

	url := "https://k-rec.360kan.com/hotrizon2/list?appid=hotrizon2&cdn_url=1&ch=AppStore&channel_id=0&ckw=0&columns=0&detection=2&direction=down&ios=11.3.1&kw=0&m2=2ae77ceb514c987b47e137895f3ebf4b&os_type=ios&sign=297E2BE4F0C4EE86A81CB04039C2CCCA&svc=3&time=1529390761906&vc=10225"
	for {
		if !GetKuaiVideoPageContent(url,db){
			break
		}
	}
}

type KuaiUrl struct {
	//Bitrate  uint32 `json:"bitrate"`
	CdnUrl   string `json:"cdn_url"`
	Height   interface{} `json:"height"`
	Size     interface{} `json:"size"`
	Width    interface{} `json:"width"`
}
type KuaiUrlList struct {
	V4g  KuaiUrl `json:"4g"`
	Wifi KuaiUrl `json:"wifi"`
}

type KuaiCnts struct {
	//FavorCnt   uint32 `json:"favorCnt"`
	//CommentCnt uint32 `json:"commentCnt"`
	ZanCnt     interface{} `json:"zanCnt"`
}

type KuaiAuthor struct {
	Qid         string `json:"qid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

type KuaiVideo struct {
	Id             string      `json:"id"`
	Title          string      `json:"title"`
	DescriptionStr string      `json:"descriptionStr"`
	Duration       interface{} `json:"duration"`
	DurationStr    string      `json:"durationStr"`
	AuthorInfo     KuaiAuthor  `json:"authorInfo"`
	PlayCnt        interface{}      `json:"playCnt"`
	Cover          string      `json:"cover"`
	CatName        string      `json:"catName"`
	Resources      KuaiUrlList `json:"resources"`
	Cnts           KuaiCnts    `json:"cnts"`
}
type KuaiData struct {
	VideoList []KuaiVideo `json:"videoList"`
}

type Kuai struct {
	Errno  string   `json:"errno"`
	Errmsg string   `json:"errmsg"`
	Data   KuaiData `json:"data"`
}

func GetKuaiVideoPageContent(apiurl string,db *gorm.DB) bool {

	requ, err := http.NewRequest("GET", apiurl, nil)
	requ.Header.Add("Host", "k-rec.360kan.com")
	requ.Header.Add("User-Agent", "LiVideoIOS/4.3.6 (iPhone; iOS 11.3.1; Scale/3.00)")

	resp, err := http.DefaultClient.Do(requ)
	if err != nil {
		logger.Debug("Proxy failed!")
		return false
	}

	recv, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return false
	}

	logger.Debug(string(recv))

	var kuai Kuai
	if err = json.Unmarshal(recv, &kuai); err != nil {
		logger.Error(err)
		return false
	}

	if kuai.Errmsg != "success" {
		logger.Error(errors.New("快视频接口返回数据异常!!!!"))
		return false
	}

	for _, v := range kuai.Data.VideoList {
		var video model.Video
		video.Provider = constant.ContentProviderKuai
		video.SourceId = v.Id
		now := time.Now()


		if err := db.Where("provider = ? and source_id = ?",video.Provider,video.SourceId).First(&video).Error ; err == nil{
			continue
		}

		tx := db.Begin()

		var category model.Category
		category.Name = v.CatName
		if err := tx.Where("name = ? and content_type = ?",category.Name,constant.MediaTypeVideo).First(&category).Error ; err == gorm.ErrRecordNotFound{
			category.ContentType = constant.MediaTypeVideo
			category.Sort = 0
			category.CreatedAt = now
			category.UpdatedAt = now
			if err := db.Save(&category).Error ;err != nil{
				logger.Error(err)
				tx.Rollback()
				return false
			}
		}

		var person model.Person
		person.Provider = constant.ContentProviderKuai
		person.SourceId = v.AuthorInfo.Qid
		if err := tx.Where("provider = ? and source_id = ?",person.Provider,person.SourceId).First(&person).Error ; err == nil{
			continue
		}
		if !strings.Contains(v.AuthorInfo.Description , "这个家伙很懒，什么也没有留下"){
			person.Description = v.AuthorInfo.Description
		}
		person.Nickname = v.AuthorInfo.Name
		person.Name = v.AuthorInfo.Name
		person.Figure = v.AuthorInfo.Avatar
		person.CreatedAt = now
		person.UpdatedAt = now
		if err = tx.Create(&person).Error ; err != nil{
			logger.Error(err)
			tx.Rollback()
			return false
		}

		video.CategoryId = category.Id

		video.PersonId = person.Id
		switch d := v.Resources.Wifi.Width.(type) {
		case string:
			width , _ := strconv.Atoi(d)
			video.Width = uint32(width)
		case int:
			video.Width = uint32(d)
		}

		switch d := v.Resources.Wifi.Height.(type) {
		case string:
			height , _ := strconv.Atoi(d)
			video.Height = uint32(height)
		case int:
			video.Height = uint32(d)
		}

		if video.Width > video.Height{
			video.Vertical = false
			video.ThumbX = v.Cover
		}else{
			video.Vertical = true
			video.ThumbY = v.Cover
		}
		video.Title = v.Title
		video.Description = v.DescriptionStr
		switch v.Duration.(type) {
		case string:
			duration , _ := strconv.Atoi(v.Duration.(string))
			video.Duration = uint32(duration)
		case int:
			video.Duration = uint32(v.Duration.(int))
		}

		switch d := v.Cnts.ZanCnt.(type) {
		case string:
			diggs , _ := strconv.Atoi(d)
			video.Diggs = uint32(diggs)
		case int:
			video.Diggs = uint32(d)
		}

		switch d := v.PlayCnt.(type) {
		case string:
			plays , _ := strconv.Atoi(d)
			video.Plays = uint32(plays)
		case int:
			video.Plays = uint32(d)
		}

		switch d := v.Resources.Wifi.Size.(type) {
		case string:
			size , _ := strconv.Atoi(d)
			video.Filesize = uint32(size)
		case int:
			video.Filesize = uint32(d)
		}

		video.Filesize = video.Filesize / 1024 / 1024

		video.Url = v.Resources.Wifi.CdnUrl
		video.Status = constant.MediaStatusReleased
		video.CreatedAt = now
		video.UpdatedAt = now

		if err = tx.Create(&video).Error ; err != nil{
			logger.Error(err)
			tx.Rollback()
			return false
		}
		tx.Commit()

	}

	return true
}
