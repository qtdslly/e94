package setting

import (
	"common/constant"
	"common/logger"
	"component/kv"
	"encoding/json"

	"github.com/jinzhu/gorm"
)

// 控制不同mms_provider的内容在不同渠道的分发
type ProviderSetting struct {
	ProviderId uint32   `json:"provider_id"`
	Channels   []string `json:"channels"`
}

func GetProviderSetting(appId uint32, db *gorm.DB) ([]*ProviderSetting, error) {
	value, err := kv.GetValueForKey(appId, 0, constant.ProviderSettingKey, false, db)
	if err != nil {
		return nil, err
	}

	if value == "" {
		return []*ProviderSetting{}, nil
	}

	var providerSet []*ProviderSetting
	err = json.Unmarshal([]byte(value), &providerSet)
	if err != nil {
		return nil, err
	}
	return providerSet, nil
}

func SetProviderSetting(appId uint32, providerSet []*ProviderSetting, db *gorm.DB) error {
	value, err := json.Marshal(providerSet)
	if err != nil {
		return err
	}

	err = kv.SetValueForKey(appId, 0, constant.ProviderSettingKey, string(value), db)
	if err != nil {
		return err
	}

	return nil
}

func SetProviderSettingByItem(appId uint32, set *ProviderSetting, db *gorm.DB) error {
	providerCache, err := GetProviderSetting(appId, db)
	if err != nil {
		logger.Error(err)
		return err
	}

	exist := false
	for index, item := range providerCache {
		if item.ProviderId == set.ProviderId {
			exist = true
			providerCache[index] = set
			break
		}
	}
	if !exist {
		providerCache = append(providerCache, set)
	}

	err = SetProviderSetting(appId, providerCache, db)
	return err
}

func GetProviderSettingByProvider(appId, providerId uint32, db *gorm.DB) *ProviderSetting {
	providerCache, err := GetProviderSetting(appId, db)
	if err != nil {
		logger.Error(err)
		return &ProviderSetting{providerId, []string{}}
	}

	for _, item := range providerCache {
		if item.ProviderId == providerId {
			return item
		}
	}

	return &ProviderSetting{providerId, []string{}}
}
