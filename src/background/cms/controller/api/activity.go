package api

import (
	"background/cms/controller/api/cache"
	"background/cms/model"
	"background/common/logger"
	"background/common/constant"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"background/cms/controller/api/model"
)

func ActivityHandler(c *gin.Context) {
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	version := c.MustGet(constant.ContextAppVersion).(*model.Version)
	appId := c.MustGet(constant.ContextAppId).(uint32)

	activity, err := LoadActivity(appId,version, db)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success,"err_msg":constant.TranslateErrCode(constant.Success),"activity":activity})
}


func LoadActivity(appId uint32,version *model.Version, db *gorm.DB) (*apimodel.AppActivity, error) {
	activity, err := cache.GetActivity(appId,version.Id,db)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	apiActivity := &apimodel.AppActivity{}
	apiActivity.Account = activity.Account
	apiActivity.Channel = activity.Channel
	apiActivity.Title = activity.Title
	apiActivity.Description = activity.Description
	apiActivity.Thumb = activity.Thumb

	return apiActivity, nil

}