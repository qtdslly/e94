package api

import (
	"background/common/logger"
	"background/newmovie/model"
	"background/common/constant"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func DiggListHandler(c *gin.Context) {
	type param struct {
		InstallationId   uint64 `form:"installation_id" binding:"required"`
		ContentType uint8  `form:"content_type" binding:"required"`
		Limit       uint32 `form:"limit" binding:"required"`
		Offset      uint32 `form:"offset" binding:"exists"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var contents []*model.ContentAction
	if err := db.Where("installation_id = ? AND content_type = ? AND content_id = ? AND action = 1", p.InstallationId, p.ContentType).Find(model.ContentAction{}).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": contents})
}

func DiggHandler(c *gin.Context) {

	type param struct {
		InstallationId   uint64 `json:"installation_id" binding:"required"`
		ContentType      uint8  `json:"content_type" binding:"required"`
		ContentId        uint32 `json:"content_id" binding:"required"`
		Disable          bool   `json:"disable"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var sVideo model.Video
	if err := db.Where("id = ?", p.ContentId).First(&sVideo).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var action model.ContentAction
	action.InstallationId = p.InstallationId
	action.ContentType = p.ContentType
	action.ContentId = p.ContentId
	action.Action = uint8(1)
	var stream model.Stream
	if err := db.Where("id = ?",action.ContentId).First(&stream).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	action.Title = stream.Title
	action.Thumb = stream.Thumb
	if p.Disable {
		if err := db.Where("installation_id = ? AND content_type = ? AND content_id = ? AND action = ?", action.InstallationId, action.ContentType, action.ContentId, action.Action).First(model.ContentAction{}).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				logger.Error(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			} else {
				if err = db.Create(&action).Error ; err != nil{
					logger.Error(err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
			}
		}
	} else {
		if err := db.Where("installation_id = ? AND content_type = ? AND content_id = ? AND action = ?", action.InstallationId, action.ContentType, action.ContentId, action.Action).First(model.ContentAction{}).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}else{
			if err := db.Where("installation_id = ? AND content_type = ? AND content_id = ? AND action = ?", action.InstallationId, action.ContentType, action.ContentId, action.Action).Delete(model.ContentAction{}).Error; err != nil {
				logger.Error(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}

