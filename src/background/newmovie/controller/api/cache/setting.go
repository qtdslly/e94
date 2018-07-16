package cache

import (
	"cms/service"
	"cms/setting"
	"common/logger"

	"github.com/jinzhu/gorm"
)

func GetFunctionSettingByChannel(appId, versionId uint32, channel string, db *gorm.DB) *setting.FunctionSetting {
	key := service.GetCacheKey("functionsetting", appId, versionId, 0, 0, channel)
	var funcSet setting.FunctionSetting
	err := service.GetCacheObject(key, &funcSet, func() (interface{}, error) {
		set, err := setting.GetFunctionSettingByChannel(appId, versionId, channel, db)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		return set, nil
	})
	if err != nil {
		logger.Error(err)
		return &funcSet
	}
	return &funcSet
}
