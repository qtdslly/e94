package setting

import (
	"common/constant"
	"common/logger"
	"component/kv"
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type RedisSetting struct {
	Type     int8   `json:"type"`
	Addr     string `json:"addr"`
	Password string `json:"password"`
}

// 获取所有的设置
func getRedisSetting(appId uint32, db *gorm.DB) ([]*RedisSetting, error) {
	value, err := kv.GetValueForKey(appId, 0, constant.RedisSettingKey, false, db)
	if err != nil {
		return nil, err
	}

	var redisSet []*RedisSetting
	err = json.Unmarshal([]byte(value), &redisSet)
	if err != nil {
		return nil, err
	}
	return redisSet, nil
}

func SetRedisSetting(appId uint32, redisSet []*RedisSetting, db *gorm.DB) error {
	value, err := json.Marshal(redisSet)
	if err != nil {
		return err
	}

	err = kv.SetValueForKey(appId, 0, constant.RedisSettingKey, string(value), db)
	if err != nil {
		return err
	}

	return nil
}

// 根据app_id获取设置
func GetRedisSettingByType(appId uint32, typ int8, db *gorm.DB) *RedisSetting {
	redisCache, err := getRedisSetting(appId, db)
	if err != nil {
		// 2017-12-21: 如果没找到则取全局设置
		if appId != 0 {
			redisCache, err = getRedisSetting(0, db)
		}

		if err != nil {
			logger.Error(err)
			return &RedisSetting{Type: typ}
		}
	}

	for _, item := range redisCache {
		if item.Type == typ {
			return item
		}
	}

	return &RedisSetting{Type: typ}
}

// 更新单条记录配置
func SetRedisSettingByItem(appId uint32, set *RedisSetting, db *gorm.DB) error {
	redisCache, err := getRedisSetting(appId, db)
	if err != nil {
		logger.Error(err)
		return err
	}

	exist := false
	for index, item := range redisCache {
		if item.Type == set.Type {
			exist = true
			redisCache[index] = set
			break
		}
	}
	if !exist {
		redisCache = append(redisCache, set)
	}

	err = SetRedisSetting(appId, redisCache, db)
	return err
}
