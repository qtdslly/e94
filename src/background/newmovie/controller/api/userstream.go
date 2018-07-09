package api


import (
	"net/http"
	"background/newmovie/model"
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"background/common/util"
)


func UserStreamAddHandler(c *gin.Context) {

	type param struct {
		InstallationId  uint64   `json:"installation_id"`
		Title           string   `json:"title"`
		Url             string   `json:"url"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var userStream model.UserStream
	userStream.Title = p.Title
	userStream.Url = p.Url
	userStream.InstallationId = p.InstallationId
	if err = db.Where("installation_id = ? and url = ?",userStream.InstallationId,userStream.Url).First(&userStream).Error ; err ==nil{
		c.JSON(http.StatusOK, gin.H{"err_code": constant.PlayurlExists, "err_msg": constant.TranslateErrCode(constant.PlayurlExists)})
	}

	userStream.Pinyin = util.TitleToPinyin(userStream.Title)
	if err = db.Create(&userStream).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}



func UserStreamUpdateHandler(c *gin.Context) {

	type param struct {
		UserStreamId    uint32   `json:"user_stream_id"`
		Title           string   `json:"title"`
		Url             string   `json:"url"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var userStream model.UserStream
	userStream.Id = p.UserStreamId
	if err = db.Where("id = ?",userStream.Id).First(&userStream).First(&userStream).Error ; err !=nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userStream.Title = p.Title
	userStream.Url = p.Url
	userStream.Pinyin = util.TitleToPinyin(userStream.Title)
	if err = db.Save(&userStream).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}


func UserStreamDeleteHandler(c *gin.Context) {

	type param struct {
		UserStreamId    uint32   `json:"user_stream_id"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	if err := db.Where("id = ?", p.UserStreamId).Delete(model.UserStream{}).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}



func UserStreamListHandler(c *gin.Context) {

	type param struct {
		InstallationId    uint64   `form:"installation_id"`
		Limit             int      `form:"limit" binding:"required"`
		Offset            int      `form:"offset" binding:"exists"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var userStreams []model.UserStream
	if err := db.Offset(p.Offset).Limit(p.Limit).Where("installation_id = ?", p.InstallationId).Find(&userStreams).Error; err != nil {
		if err != gorm.ErrRecordNotFound{
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	var count uint32
	if err := db.Model(&model.UserStream{}).Where("installation_id = ?", p.InstallationId).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	type ApiUserStream struct {
		Id              uint32     `json:"id"`
		Title           string     `json:"title"`
		Url             string     `json:"url"`
	}

	var apiUserStreams []*ApiUserStream
	for _,us := range userStreams{
		var apiUserStream ApiUserStream
		apiUserStream.Id = us.Id
		apiUserStream.Title = us.Title
		apiUserStream.Url = us.Url
		apiUserStreams = append(apiUserStreams,&apiUserStream)
	}

	var hasMore bool = true
	if len(apiUserStreams) < p.Limit{
		hasMore = false
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success,"data":apiUserStreams,"count":count,"has_more":hasMore})
}