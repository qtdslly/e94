package cms

import (
	"time"
	"background/common/constant"
	"background/common/logger"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"encoding/json"
	"background/stock/component/kv"
)

type ScriptSetting struct {
	ScriptPeriod int    `json:"script_period"` // s为单位
	Script       string `json:"script"`
	UpdatedAt    string `json:"updated_at"`
}

/*
	POST /cms/setting/script/save
	修改脚本配置
	http://localhost:2000/#!./cms/cms-setting.md
*/
func ScriptSettingSaveHandler(c *gin.Context) {
	var p ScriptSetting
	var err error

	if err = c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	// 先获取setting
	set, _ := GetScriptSetting(db)
	// 当脚本内容不一致时，自动更新时间
	if set != nil {
		if set.Script != p.Script {
			p.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		} else {
			p.UpdatedAt = set.UpdatedAt
		}
	} else {
		p.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	}

	if err := SetScriptSetting(&p, db); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}


func GetScriptSetting(db *gorm.DB) (*ScriptSetting, error) {
	value, err := kv.GetValueForKey(0, 0, constant.ScriptSettingKey, db)
	if err != nil {
		return nil, err
	}

	if value == "" {
		return nil, kv.ErrKeyNotFound
	}

	var setting ScriptSetting
	err = json.Unmarshal([]byte(value), &setting)
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func SetScriptSetting(setting *ScriptSetting, db *gorm.DB) error {
	value, err := json.Marshal(setting)
	if err != nil {
		return err
	}

	err = kv.SetValueForKey(0, 0, constant.ScriptSettingKey, string(value), db)
	if err != nil {
		return err
	}

	return nil
}