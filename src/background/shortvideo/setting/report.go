package setting

import (
	"common/constant"
	"component/kv"
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type ReportSetting struct {
	ReportWeight uint32 `json:"report_weight"` // 刷量概率, 0-100
}

func GetReportSetting(appId uint32, db *gorm.DB) (*ReportSetting, error) {
	value, err := kv.GetValueForKey(appId, 0, constant.ReportSettingKey, false, db)
	if err != nil {
		return nil, err
	}

	setting := &ReportSetting{}
	err = json.Unmarshal([]byte(value), setting)
	if err != nil {
		return nil, err
	}
	return setting, nil
}

func SetReportSetting(appId uint32, setting *ReportSetting, db *gorm.DB) error {
	value, err := json.Marshal(setting)
	if err != nil {
		return err
	}

	err = kv.SetValueForKey(appId, 0, constant.ReportSettingKey, string(value), db)
	if err != nil {
		return err
	}

	return nil
}
