package setting

import (
	"common/constant"
	"component/kv"
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type MmsProviderSetting struct {
	ProviderId []uint32 `json:"provider_id"`
}

func GetMmsProviderSetting(db *gorm.DB) (*MmsProviderSetting, error) {
	value, err := kv.GetValueForKey(0, 0, constant.MmsProviderSettingKey, false, db)
	if err != nil {
		return nil, err
	}

	if value == "" {
		set := MmsProviderSetting{ProviderId: []uint32{}}
		SetMmsProviderSetting(&set, db)
		return &set, nil
	}

	setting := &MmsProviderSetting{}
	err = json.Unmarshal([]byte(value), setting)
	if err != nil {
		return nil, err
	}
	return setting, nil
}

func SetMmsProviderSetting(setting *MmsProviderSetting, db *gorm.DB) error {
	value, err := json.Marshal(setting)
	if err != nil {
		return err
	}

	err = kv.SetValueForKey(0, 0, constant.MmsProviderSettingKey, string(value), db)
	if err != nil {
		return err
	}

	return nil
}
