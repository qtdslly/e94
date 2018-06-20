package main

import (
	"flag"
	"background/pic/config"
	"background/pic/model"
	"background/common/logger"
	"background/common/constant"
	"background/common/util"
	"github.com/tidwall/gjson"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"strings"
	"path/filepath"
	"os"
	"io"
	"fmt"
	"errors"
	"time"
)

func main(){
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

	url := "https://page.appdao.com/forward?link=16195339&style=160105&item=883716&page=1&limit=25&after=57210369&screen_w=1242&screen_h=2208&ir=0&app=1P_ElfWallpapers&v=1.4&lang=zh-Hans-CN&it=1529458959.610182&ots=3&jb=0&as=0&mobclix=0&deviceid=replaceudid&macaddr=&idv=4ADB9AA0-A063-492D-8FD6-632E835F7AF4&idvs=&ida=F7835587-9CBD-4B92-8C6E-2813811E7B5F&phonetype=iphone&model=iphone8%2C2&osn=iOS&osv=11.3.1&tz=8"
	GetBizhi(url,db)

}

func GetBizhi(url string,db *gorm.DB)bool{
	requ, err := http.NewRequest("GET", url, nil)
	requ.Header.Add("Host", "page.appdao.com")
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

	data := gjson.Get(string(recv), "data")

	if data.Exists() {
		re := data.Array()
		for _, v := range re {
			re1 := v.Array()
			for _, v1 := range re1 {
				now := time.Now()

				var picture model.Picture
				picture.SourceId = fmt.Sprint(v1.Get("link_id").Uint())
				picture.Provider = constant.PictureProivderBiZhiJinXuan
				if err = db.Where("provider = ? and source_id = ?",picture.Provider,picture.SourceId).First(&picture).Error ; err == nil {
					continue
				}

				tx := db.Begin()
				var category model.Category
				category.Name = v1.Get("link.name").String()
				if err = tx.Where("name = ?",category.Name).First(&category).Error ; err == gorm.ErrRecordNotFound{
					category.Description = v1.Get("link.description").String()
					icon := v1.Get("link.url").String()
					category.SourceIcon = icon
					if strings.Contains(icon,"jpg"){
						start := strings.Index(icon,"http")
						end := strings.Index(icon,".jpg")
						logger.Debug(start)
						logger.Debug(end)
						logger.Debug(icon)

						icon = icon[start:end]

						fileName := util.TitleToPinyin(category.Name) + ".jpg"
						category.Icon,err = DownloadFile(category.SourceIcon,config.GetStorageRoot(),fileName)
						if err != nil{
							logger.Error(err)
							tx.Rollback()
							return false
						}
					}

					category.Sort = 0
					category.CreatedAt = now
					category.UpdatedAt = now

					if err = tx.Create(&category).Error ;err != nil{
						logger.Error(err)
						tx.Rollback()
						return false
					}
				}

				var moveId uint32
				moveFile := v1.Get("movie_file").String()
				if moveFile != ""{
					var move model.Move
					move.SourceUrl = moveFile
					move.Provider = constant.PictureProivderBiZhiJinXuan
					if err = tx.Where("provider = ? and source_url = ?",move.Provider,move.SourceUrl).First(&move).Error ; err == gorm.ErrRecordNotFound {
						move.Sort = 0
						move.Title = v1.Get("title").String()
						move.Pinyin = util.TitleToPinyin(move.Title)
						move.CategoryId = category.Id
						move.Duration = float64(v1.Get("duration").Float())
						move.Click = 0
						move.Description = ""
						move.Filesize = 0
						move.Height = 0
						move.Width = 0
						move.Like = uint32(v1.Get("like").Int())
						move.OnLine = true
						move.Share = 0
						move.Show = 0

						move.SourceUrl = moveFile
						code := util.RandString(8)
						fileName := code + ".MOV"
						move.Url,err = DownloadFile(move.SourceUrl,config.GetStorageRoot(),fileName)
						if err != nil{
							logger.Error(err)
							tx.Rollback()
							return false
						}

						if move.Width > move.Height{
							move.Vertical = false
						}else{
							move.Vertical = true
						}
						if v1.Get("watermark").String() == "1"{
							move.WaterMark = true
						}else {
							move.WaterMark = false
						}
						move.CreatedAt = now
						move.UpdatedAt = now

						if err = tx.Create(&move).Error ;err != nil{
							logger.Error(err)
							tx.Rollback()
							return false
						}
					}
					moveId = move.Id
				}

				if err = tx.Where("provider = ? and source_id = ?",picture.Provider,picture.SourceId).First(&picture).Error ; err == gorm.ErrRecordNotFound {
					picture.Title = v1.Get("title").String()
					picture.SourceUrl = v1.Get("thumb_image").String()
					code := util.RandString(8)

					fileName := code + ".jpg"
					picture.Url,err = DownloadFile(picture.SourceUrl,config.GetStorageRoot(),fileName)
					if err != nil{
						logger.Error(err)
						tx.Rollback()
						return false
					}

					picture.Sort = 0
					picture.Title = v1.Get("title").String()
					picture.Pinyin = util.TitleToPinyin(picture.Title)
					picture.CategoryId = category.Id
					picture.Click = 0
					picture.Description = ""
					picture.Filesize = 0
					picture.Height = 0
					picture.Width = 0
					picture.Like = uint32(v1.Get("like").Int())
					picture.OnLine = true
					picture.Share = 0
					picture.Show = 0
					if picture.Width > picture.Height{
						picture.Vertical = false
					}else{
						picture.Vertical = true
					}
					if moveFile != "" {
						picture.IsMove = true
					}else {
						picture.IsMove = false
					}

					if v1.Get("watermark").String() == "1"{
						picture.WaterMark = true
					}else {
						picture.WaterMark = false
					}
					picture.MoveId = moveId
					picture.CreatedAt = now
					picture.UpdatedAt = now

					if err = tx.Create(&picture).Error ;err != nil{
						logger.Error(err)
						tx.Rollback()
						return false
					}
				}

				tags := strings.Split(v1.Get("tag").String(),",")
				logger.Debug(tags)
				for _, t := range tags{
					var tag model.Tag
					tag.Name = t
					if err = tx.Where("name = ?",tag.Name).First(&tag).Error ; err == gorm.ErrRecordNotFound{
						tag.Sort = 0
						tag.CreatedAt = now
						tag.UpdatedAt = now
						if err = tx.Create(&tag).Error ;err != nil{
							logger.Error(err)
							tx.Rollback()
							return false
						}
					}

					if err := tx.Exec("insert into picture_tag(picture_id,tag_id) values(?,?)", picture.Id, tag.Id).Error; err != nil {
						logger.Error(err)
						tx.Rollback()
						return false
					}

				}
				tx.Commit()
			}
		}
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
	if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		err = errors.New("invalid response content type")
		logger.Error(err)
		return "", err
	}

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
