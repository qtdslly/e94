package api

import (
	"net/http"
	"background/newmovie/model"
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func PageHandler(c *gin.Context) {

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var resourceGroups []*model.ResourceGroup
	if err = db.Order("sort asc").Where("on_line = ? and type = ?",constant.MediaStatusOnLine,constant.MediaTypeStream).Find(&resourceGroups).Error ; err != nil{
		logger.Error("query resource_group err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var streams []*model.Stream
	if err = db.Limit(15).Joins("inner join stream_group where stream.id = stream_group.stream_id and stream_group.resource_group_id = ?",resourceGroups[0].Id).Find(&streams).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var count uint32
	if err = db.Model(&model.Stream{}).Joins("inner join stream_group where stream.id = stream_group.stream_id and stream_group.resource_group_id = ?",resourceGroups[0].Id).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	type ApiStream struct {
		Id	uint32 `json:"id"`
		Title	string `json:"title"`
		Thumb	string `json:"thumb"`
	}
	type ApiPage struct {
		PageName   []string      `json:"page_name"`
		FirstPage  []*ApiStream  `json:"first_page"`
	}

	var apiPage ApiPage
	for _,resourceGroup := range resourceGroups{
		apiPage.PageName = append(apiPage.PageName,resourceGroup.Name)
	}

	for _ , stream := range streams{
		var apiStream ApiStream
		apiStream.Id = stream.Id
		apiStream.Title = stream.Title
		apiStream.Thumb = "http://www.ezhantao.com" + stream.Thumb
		apiPage.FirstPage = append(apiPage.FirstPage,&apiStream)
	}

	var hasMore bool = true
	if len(apiPage.FirstPage) < 15{
		hasMore = false
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiPage,"count":count,"has_more":hasMore})
}

