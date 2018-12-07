package api

import (
	"net/http"
	"background/newmovie/model"
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	apimodel "background/newmovie/controller/api/model"
)


func RecommendHandler(c *gin.Context) {

	type param struct {
		Limit       int `form:"limit" binding:"required"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var recommends []*model.Recommend
	if err = db.Order("created_at desc").Limit(p.Limit).Where("on_line = ?",constant.MediaStatusOnLine).Find(&recommends).Error ; err != nil{
		logger.Error("query recommend err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	type ApiRecommend struct {
		Id              uint32     `json:"id"`
		ContentType     uint32     `json:"content_type"`
		ContentId       uint32     `json:"content_id"`
		Title           string     `json:"title"`
		Focus           string     `json:"focus"`
		Thumb           string     `json:"thumb"`
	}

	var apiRecommends []*ApiRecommend
	for _,ar := range recommends{
		var apiRecommend ApiRecommend
		apiRecommend.Id = ar.Id
		apiRecommend.Title = ar.Title
		apiRecommend.ContentType = ar.ContentType
		apiRecommend.ContentId = ar.ContentId
		apiRecommend.Thumb = ar.ThumbX
		apiRecommends = append(apiRecommends,&apiRecommend)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiRecommends})
}

func VideoListHandler(c *gin.Context) {

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

	var videos []*model.Video
	if err = db.Order("publish_date desc").Offset(p.Offset).Where("on_line = ?",constant.MediaStatusOnLine).Limit(p.Limit).Find(&videos).Error ; err != nil{
		logger.Error("query movie err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	var hasMore bool = true
	if len(videos) != p.Limit{
		hasMore = false
	}

	var count uint32
	if err = db.Model(&model.Video{}).Where("on_line = ?",constant.MediaStatusOnLine).Count(&count).Error; err != nil {
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
	for _,video := range videos{
		var apiVideo ApiVideo
		apiVideo.Id = video.Id
		apiVideo.Title = video.Title
		apiVideo.Score = video.Score
		apiVideo.ThumbY = video.ThumbY
		apiVideos = append(apiVideos,&apiVideo)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiVideos,"count":count,"has_more":hasMore})
}



func VideoDetailHandler(c *gin.Context) {

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

	var video model.Video
	if err = db.Where("id = ? and on_line = ?",p.Id,constant.MediaStatusOnLine).Find(&video).Error ; err != nil{
		logger.Error("query video err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	var apiVideo *apimodel.Video
	apiVideo = apimodel.VideoFromDb(video,db)

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiVideo})
}


func VideoSearchHandler(c *gin.Context) {

	type param struct {
		Title       string `form:"title" binding:"required"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var videos []model.Video
	if err = db.Order("publish_date desc").Where("title like ? and on_line = ?","%" + p.Title + "%",constant.MediaStatusOnLine).Find(&videos).Error ; err != nil{
		logger.Error("query video err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	type ApiVideo struct {
		Id	uint32 `json:"id"`
		Title	string `json:"title"`
		Score	float64 `json:"score"`
		ThumbY	string `json:"thumb_y"`
	}

	var apiVideos []*ApiVideo
	for _,video := range videos{
		var apiVideo ApiVideo
		apiVideo.Id = video.Id
		apiVideo.Title = video.Title
		apiVideo.Score = video.Score
		apiVideo.ThumbY = video.ThumbY
		apiVideos = append(apiVideos,&apiVideo)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiVideos})
}


func VideoTopSearchHandler(c *gin.Context) {

	type param struct {
		Limit       uint32 `form:"limit" binding:"required"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var tops []model.TopSearch
	if err = db.Limit(p.Limit).Find(&tops).Error ; err != nil{
		logger.Error("query top_search err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": tops})
}
