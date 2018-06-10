package ams

import (
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"background/common/logger"
	"background/common/constant"
	"background/newmovie/model"
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
		Score string `form:"score" json:"score" ` //登录密码, password, smscode至少需要一项有值
		Actors string `form:"actors" json:"actors" ` //登录密码, password, smscode至少需要一项有值
		Directors string `form:"directors" json:"directors" ` //登录密码, password, smscode至少需要一项有值
		ThumbY string `form:"thumb_y" json:"thumb_y" ` //登录密码, password, smscode至少需要一项有值
		ThumbX string `form:"thumb_x" json:"thumb_x" ` //登录密码, password, smscode至少需要一项有值
		Url string `form:"url" json:"url" ` //登录密码, password, smscode至少需要一项有值
		PublishDate string `form:"publish_date" json:"publish_date" ` //登录密码, password, smscode至少需要一项有值
	}
	var p param
	var err error
	if err = c.Bind(&p); err != nil {
		logger.Debug("Invalid request param ", err)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var movie model.Movie
	movie.Url = p.Url
	movie.Title = p.Title
	movie.Description = p.Description
	movie.Actors = p.Actors
	movie.Directors = p.Directors
	movie.PublishDate = p.PublishDate
	movie.Score = p.Score
	movie.ThumbX = p.ThumbX
	movie.ThumbY = p.ThumbY

	now := time.Now()
	movie.CreatedAt = now
	movie.UpdatedAt = now

	// add operation log when handler return
	defer func() {
		// not log the password
		if err != nil {
			c.Set(constant.ContextError, err.Error())
		}
	}()

	if err := db.Save(&movie).Error ; err != nil{
		logger.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": movie})
}
