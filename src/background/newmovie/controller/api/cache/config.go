package cache

import (
	"background/newmovie/model"
	"background/newmovie/service"
	"background/common/logger"

	"github.com/jinzhu/gorm"
)

func GetUpgrade(db *gorm.DB) ([]*model.Upgrade, error) {
	key := service.GetCacheKey("upgrade", 0, 0, 0, 0)
	var ups []*model.Upgrade
	err := service.GetCacheObject(key, &ups, func() (interface{}, error) {
		var upgrades []*model.Upgrade
		if err := db.Where("enable=true").Order("created_at desc").Find(&upgrades).Error; err != nil {
			logger.Error(err)
			return nil, err
		}

		return upgrades, nil
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return ups, nil
}


func GetActivity(appId,versionId uint32,db *gorm.DB) (*model.Activity, error) {
	key := service.GetCacheKey("activity", 0, 0, 0, 0)
	var acs *model.Activity
	err := service.GetCacheObject(key, &acs, func() (interface{}, error) {
		var activity model.Activity
		if err := db.Where("enable=true and appId = ? and versionId = ?",appId,versionId).Order("created_at desc").First(&activity).Error; err != nil {
			logger.Error(err)
			return nil, err
		}

		return &activity, nil
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return acs, nil
}


func GetVersion(appKey, version string, db *gorm.DB) (*model.Version, error) {
	key := service.GetCacheKey("version", 0, 0, 0, 0, "os_type", 0, "app_key", appKey, "version", version)
	var v model.Version
	err := service.GetCacheObject(key, &v, func() (interface{}, error) {
		var v model.Version
		if err := db.Where("app_key=? and version = ?", appKey, version).First(&v).Error; err != nil {
			logger.Error(err)
			return nil, err
		}
		return &v, nil
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return &v, nil
}

func GetApp(appId uint32, appKey string, db *gorm.DB) (*model.App, error) {
	key := service.GetCacheKey("app", appId, 0, 0, 0, "os_type", 0, "app_key", appKey)
	var v model.App
	err := service.GetCacheObject(key, &v, func() (interface{}, error) {
		var v model.App
		var err error
		if appId > 0 {
			err = db.Where("id=?", appId).First(&v).Error
		} else {
			err = db.Where("app_key=?", appKey).First(&v).Error
		}
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		return &v, nil
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return &v, nil
}

