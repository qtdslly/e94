package cms

import (
	"background/common/constant"
	"background/common/logger"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"background/newmovie/model"
)



/*
	POST /cms/setting/script/save
	修改脚本配置
	http://localhost:2000/#!./cms/cms-setting.md
*/
func ScriptSettingSaveHandler(c *gin.Context) {

	type ScriptSetting struct {
		Script       string `form:"script"`
	}

	var p ScriptSetting
	var err error

	if err = c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	logger.Debug(p.Script)
	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var kv model.KvStore
	err = db.Where("`key`= ? and app_id = ? and version_id = ?", constant.ScriptSettingKey, 0, 0).First(&kv).Error;
	if err == gorm.ErrRecordNotFound {
		kv.AppId = 0
		kv.VersionId = 0
		kv.Key = constant.ScriptSettingKey
		kv.Value = p.Script
		if err := db.Create(&kv).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"err_code": constant.ContextError})
		}
	} else {
		if err := db.Table(model.KvStore{}.TableName()).Where("`key`= ? and app_id = ? and version_id = ?", constant.ScriptSettingKey, 0, 0).Update("value", p.Script).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"err_code": constant.ContextError})
		}
	}

	//// 先获取setting
	//set, _ := GetScriptSetting(db)
	//// 当脚本内容不一致时，自动更新时间
	//if set != nil {
	//	if set.Script != p.Script {
	//		p.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	//	} else {
	//		p.UpdatedAt = set.UpdatedAt
	//	}
	//} else {
	//	p.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	//}
	//
	//if err := SetScriptSetting(&p, db); err != nil {
	//	logger.Error(err)
	//	c.AbortWithStatus(http.StatusInternalServerError)
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}
