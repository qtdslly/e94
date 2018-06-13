package setting

import (
	"common/constant"
	"common/logger"
	"component/kv"
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type CdnSetting struct {
	Type  uint32 `json:"type"`
	CdnId uint32 `json:"cdn_id"`
}

func GetCdnSetting(appId uint32, db *gorm.DB) ([]*CdnSetting, error) {
	value, err := kv.GetValueForKey(appId, 0, constant.CdnSettingKey, false, db)
	if err != nil {
		return nil, err
	}

	if value == "" {
		return []*CdnSetting{}, nil
	}

	var cdnSet []*CdnSetting
	err = json.Unmarshal([]byte(value), &cdnSet)
	if err != nil {
		return nil, err
	}
	return cdnSet, nil
}

func SetCdnSetting(appId uint32, cdnSet []*CdnSetting, db *gorm.DB) error {
	value, err := json.Marshal(cdnSet)
	if err != nil {
		return err
	}

	err = kv.SetValueForKey(appId, 0, constant.CdnSettingKey, string(value), db)
	if err != nil {
		return err
	}

	return nil
}

// 根据type保存setting
func SetCdnSettingByItem(appId uint32, set *CdnSetting, db *gorm.DB) error {
	cdnCache, err := GetCdnSetting(appId, db)
	if err != nil {
		logger.Error(err)
		return err
	}

	exist := false
	for index, item := range cdnCache {
		if item.Type == set.Type {
			exist = true
			cdnCache[index] = set
			break
		}
	}
	if !exist {
		cdnCache = append(cdnCache, set)
	}

	err = SetCdnSetting(appId, cdnCache, db)
	return err
}

// 删除配置
func DeleteCdnSettingByApp(appId uint32, db *gorm.DB) error {
	return kv.DeleteValueForKeyByApp(appId, constant.CdnSettingKey, db)
}
