package setting

import (
	"common/constant"
	"component/kv"
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type SmsSetting struct {
	Provider   uint32 `json:"provider"`    // 短信提供商
	Template   string `json:"template"`    // 短信模板key 或者 短信模板内容
	TemplateEn string `json:"template_en"` // 英文模板
}

func GetSmsSetting(appId uint32, db *gorm.DB) (*SmsSetting, error) {
	value, err := kv.GetValueForKey(appId, 0, constant.SmsSettingKey, false, db)
	if err != nil {
		return nil, err
	}

	if value == "" {
		return &SmsSetting{}, nil
	}

	setting := &SmsSetting{}
	err = json.Unmarshal([]byte(value), setting)
	if err != nil {
		return nil, err
	}

	return setting, nil
}

func SetSmsSetting(appId uint32, setting *SmsSetting, db *gorm.DB) error {
	value, err := json.Marshal(setting)
	if err != nil {
		return err
	}

	err = kv.SetValueForKey(appId, 0, constant.SmsSettingKey, string(value), db)

	if err != nil {
		return err
	}

	return nil
}
