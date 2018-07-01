package api

import (
	"net/http"
	"background/newmovie/model"
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	apimodel "background/newmovie/controller/api/model"
	"fmt"
)

func StreamListHandler(c *gin.Context) {

	type param struct {
		Limit       int `form:"limit" binding:"required"`
		Offset      int `form:"offset" binding:"exists"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var streams []*model.Stream
	if err = db.Order("sort asc").Offset(p.Offset).Where("on_line = ?",constant.MediaStatusOnLine).Limit(p.Limit).Find(&streams).Error ; err != nil{
		logger.Error("query movie err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	var hasMore bool = true
	if len(streams) != p.Limit{
		hasMore = false
	}

	var count uint32
	if err = db.Model(&model.Stream{}).Where("on_line = ?",constant.MediaStatusOnLine).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	type ApiVideo struct {
		Id	uint32 `json:"id"`
		Title	string `json:"title"`
		Score	float64 `json:"score"`
		ThumbY	string `json:"thumb_y"`
	}

	var apiVideos []*ApiVideo
	for _,stream := range streams{
		var apiVideo ApiVideo
		apiVideo.Id = stream.Id
		apiVideo.Title = stream.Title
		apiVideo.ThumbY = "http://www.ezhantao.com/thumb/stream/" + fmt.Sprint(stream.Id) + ".jpg"
		apiVideos = append(apiVideos,&apiVideo)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiVideos,"count":count,"has_more":hasMore})
}



func StreamDetailHandler(c *gin.Context) {

	type param struct {
		Id       uint32 `form:"id" binding:"required"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var stream model.Stream
	if err = db.Where("id = ? and on_line = ?",p.Id,constant.MediaStatusOnLine).Find(&stream).Error ; err != nil{
		logger.Error("query video err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	var apiVideo *apimodel.Video
	apiVideo.Id = stream.Id
	apiVideo.ThumbX = "http://www.ezhantao.com/thumb/stream/" + fmt.Sprint(stream.Id) + ".jpg"
	apiVideo.Title = stream.Title
	var playUrls []*model.PlayUrl
	if err := db.Where("content_type = 4 and content_id = ?",stream.Id).Find(&playUrls).Error ; err != nil{
		logger.Error("query play_url err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	for _,playUrl := range playUrls{
		var pUrl apimodel.PlayUrl
		pUrl.Id = playUrl.Id
		pUrl.Provider = playUrl.Provider
		pUrl.IsPlay = true
		pUrl.Url = playUrl.Url
		apiVideo.Urls = append(apiVideo.Urls,&pUrl)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiVideo})
}

