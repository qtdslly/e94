package api

import (
	"common/constant"
	"common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"

	cmod "background/cms/model"
	"background/newmovie/model"
)

func PlayUrlListHandler(c *gin.Context) {

	type param struct {
		ContentType     uint32 `form:"content_type" json:"content_type"`
		ContentId       uint32 `form:"content_id" json:"content_id"`

		//Limit  int `form:"limit" binding:"required"`
		//Offset int `form:"offset" binding:"exists"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var playUrls []cmod.PlayUrl
	if err = db.Order("sort asc").Where("content_type = ? and content_id = ?", p.ContentType,p.ContentId).Find(&playUrls).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var count uint32
	if err = db.Model(&cmod.PlayUrl{}).Where("content_type = ? and content_id = ?", p.ContentType,p.ContentId).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": playUrls, "count":count})
}

func PlayUrlAddHandler(c *gin.Context) {
	type param struct {
		Provider    uint32         `form:"provider" json:"provider"`
		ContentType uint8          `form:"content_type" json:"content_type"`
		ContentId   uint32         `form:"content_id" json:"content_id"`
		Title       string         `form:"title" json:"title"`
		Url         string         `form:"url" json:"url"`
		PageUrl     string         `form:"page_url" json:"page_url"`
		OnLine      bool           `form:"on_line" json:"on_line"` // 链接播放不了，临时禁止
		Quality     uint8          `form:"quality" json:"quality"`
		Sort        uint32         `form:"sort" json:"sort"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var playUrl cmod.PlayUrl
	if err = db.Where("content_type = ? and content_id = ? and url = ?", p.ContentType, p.ContentId,p.Url).First(&playUrl).Error; err == nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	playUrl.Title = p.Title
	playUrl.Provider = p.Provider
	playUrl.ContentType = p.ContentType
	playUrl.ContentId = p.ContentId
	playUrl.Url = p.Url
	playUrl.PageUrl = p.PageUrl
	playUrl.OnLine = p.OnLine
	playUrl.Quality = p.Quality
	playUrl.Sort = p.Sort

	if err := db.Create(&playUrl).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}


func PlayUrlDeleteHandler(c *gin.Context) {
	type param struct {
		Id          uint32         `form:"id" json:"id"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	if err := db.Delete(model.PlayUrl{},"id = ?",p.Id).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}


