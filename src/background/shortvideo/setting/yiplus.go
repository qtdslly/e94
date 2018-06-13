package setting

import (
	"common/constant"
	"component/kv"
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type YiPlusSetting struct {
	AccessKey       string `json:"access_key"`        //Yi+访问key
	SecretKey       string `json:"secret_key"`        //Yi+加密key
	SignKey         string `json:"sign_key"`          //我们接口签名key
	Domain          string `json:"domain"`            //Yi+域名
	MinDuration     uint32 `json:"min_duration"`      //直播识别最小连续出现次数
	MinInterval     uint32 `json:"min_interval"`      //直播识别最小时间间隔，间隔时间内都算连续
	StreamResultTTL uint32 `json:"stream_result_ttl"` //直播识别结果保存时间，单位s
}

func GetYiPlusSetting(db *gorm.DB) (*YiPlusSetting, error) {
	value, err := kv.GetValueForKey(0, 0, constant.YiPlusSettingKey, false, db)
	if err != nil {
		return nil, err
	}

	if value == "" {
		return &YiPlusSetting{}, nil
	}

	var yiPlusSetting *YiPlusSetting
	err = json.Unmarshal([]byte(value), &yiPlusSetting)
	if err != nil {
		return nil, err
	}
	return yiPlusSetting, nil
}

func SetYiPlusSetting(yiPlusSetting *YiPlusSetting, db *gorm.DB) error {
	value, err := json.Marshal(yiPlusSetting)
	if err != nil {
		return err
	}

	err = kv.SetValueForKey(0, 0, constant.YiPlusSettingKey, string(value), db)
	if err != nil {
		return err
	}

	return nil
}
