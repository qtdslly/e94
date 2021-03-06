package main

import (
	"background/shortvideo/config"
	"background/common/logger"
	"background/common/constant"
	"background/common/util"
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
	"fmt"
	"path/filepath"
	"os"
	"io"
)

func main() {
	logger.SetLevel(config.GetLoggerLevel())


	flag.Parse()

	logger.SetLevel(config.GetLoggerLevel())

	db, err := gorm.Open("mysql", "root:hahawap@tcp(localhost:3306)/shortvideo?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}

	db.LogMode(true)
	model.InitModel(db)
	logger.SetLevel(config.GetLoggerLevel())

	url := "https://k-rec.360kan.com/hotrizon2/list?appid=hotrizon2&cdn_url=1&ch=AppStore&channel_id=0&ckw=0&columns=0&detection=2&direction=down&ios=11.3.1&kw=0&m2=2ae77ceb514c987b47e137895f3ebf4b&os_type=ios&sign=297E2BE4F0C4EE86A81CB04039C2CCCA&svc=3&time=1529390761906&vc=10225"
	for {
		//if !GetKuaiVideoPageContent(url,db){
		//	break
		//}
		var p = time.Now()
		if fmt.Sprintf("%02d%02d",p.Hour(),p.Minute()) == "1101"{
			break
		}
		GetKuaiVideoPageContent(url,db)
	}
}

type KuaiUrl struct {
	//Bitrate  uint32 `json:"bitrate"`
	CdnUrl   string `json:"cdn_url"`
	Height   interface{} `json:"height"`
	Size     uint32 `json:"size"`
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
			tx.Rollback()
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
		var value string
		value = GetValue(v.Resources.Wifi.Width)
		if value != ""{
			width,_ := strconv.Atoi(value)
			video.Width = uint32(width)
		}

		value = GetValue(v.Resources.Wifi.Height)
		if value != ""{
			height,_ := strconv.Atoi(value)
			video.Height = uint32(height)
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
		value = GetValue(v.Duration)
		if value != ""{
			duration,_ := strconv.Atoi(value)
			video.Duration = uint32(duration)
		}

		value = GetValue(v.Cnts.ZanCnt)
		if value != ""{
			diggs,_ := strconv.Atoi(value)
			video.Diggs = uint32(diggs)
		}

		value = GetValue(v.PlayCnt)
		if value != ""{
			plays,_ := strconv.Atoi(value)
			video.Plays = uint32(plays)
		}

		video.Filesize = v.Resources.Wifi.Size

		video.Url = v.Resources.Wifi.CdnUrl

		code := util.RandString(6)
		fileName := fmt.Sprintf("%04d%02d%02d%02d%02d%02d",now.Year(),now.Month(),now.Day(),now.Hour(),now.Minute(),now.Second()) + code + ".mp4"
		video.FileName, err = DownloadFile(video.Url,config.GetStorageRoot(),fileName)
		if err != nil{
			logger.Error(err)
			tx.Rollback()
			return false
		}
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





func DownloadFile(requrl string,rootPath , filename string) (string, error) {

	//no timeout
	client := http.Client{}

	resp, err := client.Get(requrl)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	// close body read before return
	defer resp.Body.Close()

	// should not save html content as file
	//if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
	//	err = errors.New("invalid response content type")
	//	logger.Error(err)
	//	return "", err
	//}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Fail to download: [%s]", requrl))
		logger.Error(err)
		return "", err
	}

	//now := time.Now()
	//middlePath := fmt.Sprint(now.Year()) + "/" + fmt.Sprint(now.Month()) + "/" + fmt.Sprint(now.Day()) + "/" + fmt.Sprint(now.Hour()) + "/" + fmt.Sprint(now.Minute())

	relDir := time.Now().Format("/2006/01/02/15/04")
	relPath := filepath.Join(relDir, filename)

	if err := os.MkdirAll(filepath.Join(rootPath, relDir), 0755); err != nil {
		logger.Error(err)
		return "", err
	}

	tmpPath := filepath.Join(rootPath,relPath)
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	_ ,err = io.Copy(tmpFile, resp.Body)
	tmpFile.Close()
	if err != nil {
		logger.Error(err)
		return "", err
	}

	logger.Debug(tmpPath)
	return tmpPath, nil
}


func GetValue(general interface{})(string) {
	switch general.(type) {
	case uint32:
		newInt, ok := general.(uint32)
		if ok == false{
			return ""
		}
		return fmt.Sprint(newInt)
	case uint64:
		newInt, ok := general.(uint64)
		if ok == false{
			return ""
		}
		return fmt.Sprint(newInt)
	case int :
		newInt, ok := general.(int)
		if ok == false{
			return ""
		}
		return fmt.Sprint(newInt)
	case float32:
		newFloat32, ok := general.(float32)
		if ok == false{
			return ""
		}
		return fmt.Sprint(newFloat32)

	case float64:
		newFloat64, ok := general.(float64)
		if ok == false{
			return ""
		}
		return fmt.Sprint(newFloat64)

	default:
		return ""
	}
	return ""
}