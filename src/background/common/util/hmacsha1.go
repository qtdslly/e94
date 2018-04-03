package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

func HmacSha1(input string, secretKey string) string {
	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(input))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
