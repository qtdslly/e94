package sms

import (
	"common/logger"
	"testing"
)

func TestSendVerificationCode(t *testing.T) {
	code := "123456"
	var datas []string
	datas = append(datas, code) // smscode
	datas = append(datas, "10") // limited activation time

	response := SendSMSYtx("3083", datas, "15623307163")
	if response.Status != Success {
		logger.Error(response)
	}

	logger.Info("sms sent successfully")
}
