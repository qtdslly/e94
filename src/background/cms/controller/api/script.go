package api

import (
	"background/common/constant"
	"background/common/logger"
	"background/cms/model"
	"net/http"
	//"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

/*
	GET /cms/v1.0/script
	获取脚本内容
*/
func ScriptHandler(c *gin.Context) {
	type param struct {
		LastUpdate int64 `form:"last_update"`
	}

	var p param
	var err error

	if err = c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	//lastUpdate := time.Unix(p.LastUpdate, 0)

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var setting model.KvStore
	if err := db.Where("app_id = 0 and version_id = 0 and `key` = ?",constant.ScriptSettingKey).First(&setting).Error ; err != nil{
		logger.Error(err)
		return
	}
	// 只有当脚本有内容且更新时间在客户端更新时间之后，才返回脚本内容
	//if setting != nil {
	//	updatedAt, _ := time.ParseInLocation("2006-01-02 15:04:05", setting.UpdatedAt, time.Local)
	//	if updatedAt.After(lastUpdate) {
	//		c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": setting.Value})
	//		return
	//	}
	//}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": setting.Value})

	//c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}
