package sms

import (
	"bytes"
	"common/logger"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	apiServer    = "https://app.cloopen.com:8883"
	apiVersion   = "2013-12-26"
	accountSid   = "aaf98f894700d34e014713f76b8e041b"
	accountToken = "cd1a568bbe0f47e89fee9bb274af4999"
	appId        = "aaf98f89471ea2c10147240e71470178"
	Success      = "000000"
)

type smsYtxRequest struct {
	To         string   `json:"to"`
	AppId      string   `json:"appId"`
	TemplateId string   `json:"templateId"`
	Datas      []string `json:"datas"`
}

type SmsYtxResponse struct {
	Status string `json:"statusCode"` // "000000" is success
}

/*
	MD5 hash
*/
func md5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))
}

/*
	This method allows to call alidayu to send SMS
*/
func SendSMSYtx(templateId string, datas []string, mobile string) *SmsYtxResponse {
	var sig string
	now := time.Now().Format("20060102150405")
	sig = md5Hash(accountSid + accountToken + now)

	// generate api url
	apiUrl := apiServer + "/" + apiVersion + "/Accounts/" + accountSid + "/SMS/TemplateSMS?sig=" + sig

	// generate authorization token
	authSrc := accountSid + ":" + now
	auth := base64.StdEncoding.EncodeToString([]byte(authSrc))

	request := smsYtxRequest{}
	request.To = mobile
	request.AppId = appId
	request.TemplateId = templateId
	request.Datas = datas

	// marshal data
	b, _ := json.Marshal(request)
	body := bytes.NewReader(b)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", apiUrl, body)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	var response SmsYtxResponse
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
