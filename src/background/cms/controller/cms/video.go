package cms

import (
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"background/common/logger"
	"background/common/constant"
	"background/cms/model"
	"background/common/util"
	"net/http"
	"time"
)

/*
	POST /movie/save
	保存视频信息
	@Author:LLY
	http://localhost:2000/#!./cms-movie.md
*/
func MovieSaveHandler(c *gin.Context) {
	type param struct {
		Title  string `form:"title"  json:"title" binding:"required"`  //username 或 mobile 或 email
		Description string `form:"description" json:"description" binding:"required"` //登录密码, password, smscode至少需要一项有值
		Score float64 `form:"score" json:"score" ` //登录密码, password, smscode至少需要一项有值
		Actors string `form:"actors" json:"actors" ` //登录密码, password, smscode至少需要一项有值
		Directors string `form:"directors" json:"directors" ` //登录密码, password, smscode至少需要一项有值
		ThumbY string `form:"thumb_y" json:"thumb_y" ` //登录密码, password, smscode至少需要一项有值
		ThumbX string `form:"thumb_x" json:"thumb_x" ` //登录密码, password, smscode至少需要一项有值
		Url string `form:"url" json:"url" ` //登录密码, password, smscode至少需要一项有值
		PublishDate string `form:"publish_date" json:"publish_date" ` //登录密码, password, smscode至少需要一项有值
		Year uint32 `form:"year" json:"year" ` //登录密码, password, smscode至少需要一项有值
		Tags string `form:"tags" json:"tags" ` //登录密码, password, smscode至少需要一项有值
		Language string `form:"language" json:"language" ` //登录密码, password, smscode至少需要一项有值
		Country string `form:"country" json:"country" ` //登录密码, password, smscode至少需要一项有值
		Duration uint32 `form:"duration" json:"duration" ` //登录密码, password, smscode至少需要一项有值

	}
	var p param
	var err error
	if err = c.Bind(&p); err != nil {
		logger.Debug("Invalid request param ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	now := time.Now()

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var video model.Video
	video.Title = p.Title
	if err := db.Where("title = ?",video.Title).First(&video).Error ; err == gorm.ErrRecordNotFound{
		video.Description = p.Description
		video.Actors = p.Actors
		video.Directors = p.Directors
		video.PublishDate = p.PublishDate
		video.Score = p.Score
		video.ThumbX = p.ThumbX
		video.ThumbY = p.ThumbY
		video.Country = p.Country
		video.Language = p.Language
		video.Tags = p.Tags
		video.Pinyin = util.TitleToPinyin(video.Title)
		video.Year = p.Year
		video.TotalEpisode = 1
		video.OnLine = constant.MediaStatusOnLine

		video.CreatedAt = now
		video.UpdatedAt = now

		if err := db.Create(&video).Error ; err != nil{
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	var episode model.Episode
	episode.Title = p.Title
	if err := db.Where("title = ?",episode.Title).First(&episode).Error ; err == gorm.ErrRecordNotFound {
		episode.VideoId = video.Id
		episode.Pinyin = video.Pinyin
		episode.Score = video.Score
		episode.Duration = p.Duration * 60
		episode.Description = p.Description
		episode.ThumbY = p.ThumbY
		episode.PublishDate = p.PublishDate
		episode.CreatedAt = now
		episode.UpdatedAt = now

		if err := db.Create(&episode).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	var playUrl model.PlayUrl
	playUrl.ContentType = constant.MediaTypeEpisode
	playUrl.ContentId = episode.Id
	playUrl.Provider = constant.ContentProviderSystem

	if err := db.Where("content_type = ? and content_id = ? and provider = ?",playUrl.ContentType,playUrl.ContentId,playUrl.Provider).First(&playUrl).Error ; err == gorm.ErrRecordNotFound {
		playUrl.Url = p.Url
		playUrl.OnLine = true
		playUrl.Title = p.Title
		if err := db.Save(&playUrl).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": video})
}


