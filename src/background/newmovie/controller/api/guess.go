package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"background/common/logger"
	"background/common/constant"
	"background/newmovie/model"
	"background/newmovie/service"
	"strings"
)

func GuessListHandler(c *gin.Context) {
	type param struct {
		StreamId uint8  `form:"id" binding:"required"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var stream model.Stream
	if err := db.Where("id = ?",p.StreamId).First(&stream).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	areaTitle := ""
	logger.Debug(len(service.AREA))
	for _,area := range service.AREA{
		logger.Debug(stream.Title,"|",area)
		if strings.Contains(stream.Title,area){
			areaTitle = area
			break
		}
	}
	logger.Debug(areaTitle)
	if areaTitle == ""{
		areaTitle = stream.Title
	}
	logger.Debug(areaTitle)

	var streams []model.Stream
	if err := db.Limit(6).Where("id <> ? and on_line = 1 and title like ?",p.StreamId,"%" + areaTitle + "%").Find(&streams).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	
	var streams1 []model.Stream
	if len(streams) < 6{
		if err := db.Limit(6 - len(streams)).Where("id <> ? and on_line = 1 and category = ?",p.StreamId,stream.Category).Find(&streams1).Error ; err != nil{
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	
	for _ , s := range streams1{
		streams = append(streams,s)
	}

	type ApiStream struct {
		Id    uint32 `json:"id"`
		Title string `json:"title"`
		Thumb string `json:"thumb"`
	}

	var apiStreams []*ApiStream
	for _, stream := range streams {
		var apiStream ApiStream
		apiStream.Id = stream.Id
		apiStream.Title = stream.Title
		apiStream.Thumb = stream.Thumb
		apiStreams = append(apiStreams, &apiStream)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiStreams,"count":len(apiStreams),"has_more":false})
}
