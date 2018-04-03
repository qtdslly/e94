package sms

import (
	"bytes"
	"common/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	apikey = "5d9e73abb195954e0a4f20c1b75f8cf1"
	// englishTemplate  = "【China TV】Your verification code : %s"
	// chineseTemplate  = "【China TV】验证码 : %s"
	// chinaCountryCode = "+86"
)

type smsYpRequest struct {
	ApiKey string `form:"apikey"`
	Mobile string `form:"mobile"`
	Text   string `form:"text"`
}

type SmsYpResponse struct {
	Code   int     `json:"code"`
	Msg    string  `json:"msg"`
	Count  int     `json:"count"`  // 成功发送的短信计费条数
	Fee    float32 `json:"fee"`    // 扣费条数，70个字一条，超出70个字时按每67字一条计
	Unit   string  `json:"unit"`   // 计费单位
	Mobile string  `json:"mobile"` // 发送手机号
	Sid    int     `json:"sid"`    // 短信ID
}

/*
	This method allows to call alidayu to send SMS
*/
func SendSMSYp(template, countryCode, mobile, code string) *SmsYpResponse {
	// generate api url
	apiUrl := "https://sms.yunpian.com/v2/sms/single_send.json"

	content := fmt.Sprintf(template, code)

	// combine the country code and mobile number
	var combinedMobile string
	if strings.HasPrefix(mobile, "0") {
		combinedMobile = countryCode + strings.TrimPrefix(mobile, "0")
	} else {
		combinedMobile = countryCode + mobile
	}

	payload := fmt.Sprintf("apikey=%s&mobile=%s&text=%s", apikey, url.QueryEscape(combinedMobile), content)
	logger.Debug(payload)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", apiUrl, bytes.NewReader([]byte(payload)))
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var response SmsYpResponse
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return nil
	}

	defer resp.Body.Close()
	rb, err := ioutil.ReadAll(resp.Body)
	logger.Debug("sms response :" + string(rb))
	if err != nil {
		logger.Error(err)
		return nil
	}

	if err = json.Unmarshal(rb, &response); err != nil {
		logger.Error(err)
		return nil
	}
	return &response
}
