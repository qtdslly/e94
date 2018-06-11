package main

import (
	"background/common/logger"
	"background/shortvideo/config"
	"background/shortvideo/model"
	"background/common/constant"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"path"
	"flag"
	"strconv"
	"strings"
)

var ShortVideoGroupId uint32 = 1

func main(){
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

	AddShortVideo(db)
}

func AddShortVideo(db *gorm.DB){
	tx := db.Begin()
	var err error

	var thirdVideos []model.ThirdVideo

	if err = tx.Where("created_at not like '2018-06-11%'").Find(&thirdVideos).Error ; err != nil{
		logger.Error(err)
		tx.Rollback()

		return
	}
	for _,thirdVideo := range thirdVideos{
		dur,_ := strconv.Atoi(thirdVideo.Duration)
		duration := uint32(dur / 1000)

		var video model.ShortVideo
		if thirdVideo.Provider == "douyin"{
			video.IsVerticalScreen = true
			video.Provider = uint32(constant.ContentProviderDouYin)
			video.Watermark = thirdVideo.HasWaterMark
		}else{
			video.IsVerticalScreen = false
			video.Provider = uint32(constant.ContentProviderPear)
			video.Watermark = true
		}

		video.SourceId = thirdVideo.ThirdVideoId
		if err := db.Where("provider = ? and source_id = ?",video.SourceId).First(&video).Error ; err == nil{
			continue
		}
		video.Status = uint32(constant.MediaStatusReleased)
		now := time.Now()
		video.CreatedAt = now
		video.UpdatedAt = now
		video.ReleasedAt = &now
		video.Duration = duration
		video.ReleasedAt = &now
		video.ThumbY = thirdVideo.ThumbX
		video.Duration = duration
		video.Width = uint32(thirdVideo.Width)
		video.Height = uint32(thirdVideo.Height)
		video.Filesize = uint32(thirdVideo.Filesize)
		video.Url = "/episode/2018/06/08/" + path.Base(thirdVideo.FileName)
		video.Diggs = thirdVideo.DiggCount
		video.Plays = thirdVideo.PlayCount
		video.Province = thirdVideo.Province
		video.City = thirdVideo.City
		video.District = thirdVideo.District
		video.Address = thirdVideo.Address
		video.Country = thirdVideo.Country
		video.Latitude = thirdVideo.Latitude
		video.Longitude = thirdVideo.Longitude

		//video.CommentCount = thirdVideo.CommentCount
		//video.ShareCount = thirdVideo.ShareCount
		//video.ShareTitle = strings.Replace(thirdVideo.ShareTitle,"抖音","",-1)
		//video.ShareDescription = ""
		//video.AuthorThumb = thirdVideo.AuthorThumb
		//video.Birthday = thirdVideo.Birthday
		//video.NickName = thirdVideo.NickName
		//video.ThirdAuthorId = thirdVideo.ThirdAuthorId
		//video.ThirdVideoId = thirdVideo.ThirdVideoId


		if err = tx.Create(&video).Error ; err != nil{
			logger.Error(err)
			tx.Rollback()

			return
		}

		if err := db.Exec("insert into short_video_group(short_video_id,resource_group_id) values(?,?)",video.Id,ShortVideoGroupId).Error ; err != nil{
			logger.Error(err)
			tx.Rollback()

			return
		}

		var person model.Person
		person.Nickname = thirdVideo.NickName
		if err = db.Where("nickname = ? and provider_id = ?",person.Nickname,constant.ContentProviderDouYin).First(&person).Error ; err == gorm.ErrRecordNotFound {
			//if len(thirdVideo.Birthday) > 0 {
			//	thirdVideo.Birthday = thirdVideo.Birthday + " 00:00:00"
			//	timeLayout := "2006-01-02 15:04:05"                             //转化所需模板
			//	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
			//	theTime, _ := time.ParseInLocation(timeLayout, thirdVideo.Birthday, loc)
			//	person.Birthday = &theTime
			//}
			person.Country = "中国"
			person.Nickname = thirdVideo.NickName
			person.ProviderId = video.Provider
			person.Figure = thirdVideo.AuthorThumb
			person.CreatedAt = now
			person.UpdatedAt = now
			person.SyncedAt = &now

			if err = tx.Create(&person).Error; err != nil {
				logger.Error(err)
				tx.Rollback()
				return
			}
		}

		var personMedia model.PersonMedia
		personMedia.ContentType = constant.MediaTypeShortVideo
		personMedia.ContentId = video.Id
		personMedia.PersonId = person.Id

		if err = db.Where("content_type = ? and content_id = ? and person_id = ?",personMedia.ContentType,personMedia.ContentId,personMedia.PersonId).First(&personMedia).Error ; err == gorm.ErrRecordNotFound {
			personMedia.CreatedAt = now
			personMedia.UpdatedAt = now
			if err = tx.Create(&personMedia).Error; err != nil {
				logger.Error(err)
				tx.Rollback()
				return
			}
		}

		var apps []model.App
		if err := db.Find(&apps).Error ; err != nil{
			logger.Error(err)
			tx.Rollback()

			return
		}

		for _,app := range apps{
			var param model.ResourceParam
			param.ContentType = constant.MediaTypeShortVideo
			param.ContentId = video.Id
			param.AppId = app.Id

			if err = db.Where("content_type = ? and content_id = ? and app_id = ?",param.ContentType,param.ContentId,param.AppId).First(&param).Error ; err == gorm.ErrRecordNotFound {
				param.EnableHls = true
				param.Online = true
				param.CreatedAt = now
				param.UpdatedAt = now
				param.ExpiredAt = &now
				param.ReleasedAt = &now
				if err = tx.Create(&param).Error; err != nil {
					logger.Error(err)
					tx.Rollback()
					return
				}
			}
		}


		if thirdVideo.Provider == "pear"{
			ts := strings.Split(thirdVideo.Tag,",")
			for _ , videoType := range ts{

				var tag model.Tag
				tag.Name = videoType
				if err := tx.Where("name = ? and property_id = ?",tag.Name,pro["类型"]).First(&tag).Error ; err == gorm.ErrRecordNotFound{
					//tag.PropertyId = pro["类型"]
					//tag.Sort = 0
					//now := time.Now()
					//tag.CreatedAt = now
					//tag.UpdatedAt = now
					//if err = tx.Create(&tag).Error ; err != nil{
					//	logger.Error(err)
					//	tx.Rollback()
					//	return
					//}

					continue
				}

				var entityTag model.EntityTag
				entityTag.Entity = "video"
				entityTag.TagId = tag.Id
				entityTag.EntityId = video.Id

				if err = db.Where("entity = ? and tag_id = ? and entity_id = ?",entityTag.Entity,entityTag.TagId,entityTag.EntityId).First(&entityTag).Error ; err == gorm.ErrRecordNotFound{
					entityTag.PropertyId = tag.PropertyId
					entityTag.Sort = 0
					entityTag.CreatedAt = now
					entityTag.UpdatedAt = now
					if err = tx.Create(&entityTag).Error ; err != nil{
						logger.Error(err)
						tx.Rollback()
						return
					}
				}
			}
		}

	}
	tx.Commit()
}


