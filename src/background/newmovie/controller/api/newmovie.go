package api

import (
	"net/http"
	"background/newmovie/model"
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"background/newmovie/service"
)

func NewMovieListHandler(c *gin.Context) {

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

	var movies []*model.Movie
	if err = db.Order("publish_date desc").Offset(p.Offset).Where("provider <> 'youku'").Limit(p.Limit).Find(&movies).Error ; err != nil{
		logger.Error("query movie err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	var hasMore bool = true
	if len(movies) != p.Limit{
		hasMore = false
	}

	var count uint32
	if err = db.Model(&model.Movie{}).Where("provider <> 'youku'").Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//type ApiMovie struct {
	//	Id	int `json:"id"`
	//	Title	string `json:"title"`
	//	Score	string `json:"score"`
	//	ThumbY	string `json:"thumb_y"`
	//}

	//var apiMovies []*ApiMovie
	//for _,movie := range movies{
	//	var apiMovie ApiMovie
	//	apiMovie.Id = movie.Id
	//	apiMovie.Title = movie.Title
	//	apiMovie.Score = movie.Score
	//	apiMovie.ThumbY = movie.ThumbY
	//	apiMovies = append(apiMovies,apiMovie)
	//}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": movies,"count":count,"has_more":hasMore})
}



func NewMovieHandler(c *gin.Context) {

	type param struct {
		Id       int `form:"id" binding:"required"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var movie model.Movie
	if err = db.Where("id = ? and provider <> 'youku'",p.Id).Find(&movie).Error ; err != nil{
		logger.Error("query movie err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	var set model.KvStore
	if err := db.Where("`key` = 'script_setting_key'").First(&set).Error ; err != nil{
		logger.Error(err)
		return
	}

	movie.Url = service.GetRealUrl(movie.Provider,movie.Url,set.Value)

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": movie})
}


func NewMovieSearchHandler(c *gin.Context) {

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

	var movies []model.Movie
	if err = db.Order("publish_date desc").Where("title like ? and provider <> 'youku'","%" + p.Title + "%").Find(&movies).Error ; err != nil{
		logger.Error("query movie err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	var set model.KvStore
	if err := db.Where("`key` = 'script_setting_key'").First(&set).Error ; err != nil{
		logger.Error(err)
		return
	}

	var apiMovies []model.Movie
	for _,movie := range movies{
		movie.Url = service.GetRealUrl(movie.Provider,movie.Url,set.Value)
		if movie.Url != ""{
			apiMovies = append(apiMovies,movie)
		}
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiMovies})
}


func NewMovieTopSearchHandler(c *gin.Context) {

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
		logger.Error("query movie err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": tops})
}
