package api


import (
	"net/http"
	"background/cms/model"
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"time"
)


func NotifcationHandler(c *gin.Context) {
	type param struct {
		LastTime string `form:"last_time"`
	}

	var p param
	var err error

	if err = c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	startAt, err := time.ParseInLocation("2006-01-02 15:04:05", p.LastTime, time.Local)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var notifications []*model.Notification
	if err = db.Order("created_at desc").Limit(3).Where("on_line = ? and created_at > ?",constant.MediaStatusOnLine,startAt).Find(&notifications).Error ; err != nil{
		logger.Error("query recommend err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	type ApiNoti struct {
		Id              uint32     `json:"id"`
		Title           string     `json:"title"`
		Description     string     `json:"description"`
		ContentType     uint32     `json:"content_type"`
		ContentId       uint32     `json:"content_id"`
		Thumb           string     `json:"thumb"`
	}

	var apiNotis []*ApiNoti
	for _,ar := range notifications{
		var apiNoti ApiNoti
		apiNoti.Id = ar.Id
		apiNoti.Title = ar.Title
		apiNoti.ContentType = ar.ContentType
		apiNoti.ContentId = ar.ContentId
		apiNoti.Description = ar.Description
		apiNoti.Thumb = ar.Thumb
		apiNotis = append(apiNotis,&apiNoti)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiNotis})
}