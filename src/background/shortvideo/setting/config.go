package setting

import (
	"background/common/constant"
	"background/common/logger"
	"background/component/kv"
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type ConfigSetting struct {
	CmnBindAddr        string `json:"cmn_bind_addr"`
	ResBindAddr        string `json:"res_bind_addr"`
	CacheRedisAddr     string `json:"cache_redis_addr"`
	CacheRedisPassword string `json:"cache_redis_password"`
}

func GetConfigSetting(db *gorm.DB) (*ConfigSetting, error) {
	value, err := kv.GetValueForKey(0, 0, constant.ConfigSettingKey, db)
	if err != nil {
		return nil, err
	}

	setting := &ConfigSetting{}
	err = json.Unmarshal([]byte(value), setting)
	if err != nil {
		return nil, err
	}
	return setting, nil
}

func GetCmnBindAddr(db *gorm.DB) (string, error) {
	confSet, err := GetConfigSetting(db)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	return confSet.CmnBindAddr, nil
}

func GetResBindAddr(db *gorm.DB) (string, error) {
	confSet, err := GetConfigSetting(db)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return confSet.ResBindAddr, nil
}

func GetCacheRedis(db *gorm.DB) (string, string, error) {
	confSet, err := GetConfigSetting(db)
	if err != nil {
		logger.Error(err)
		return "", "", err
	}
	return confSet.CacheRedisAddr, confSet.CacheRedisPassword, nil
}
