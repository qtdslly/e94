package util

import (
	"background/common/logger"
	"fmt"
	"net/url"
)

func GetEncodedUrl(str string) string {
	if str == "" {
		return str
	}
	u, err := url.Parse(str)
	if err != nil {
		logger.Error(err)
		return str
	}
	result := ""
	if len(u.Query()) > 0 {
		result = fmt.Sprintf(`%s://%s%s?%s`, u.Scheme, u.Host, u.EscapedPath(), u.Query().Encode())
	} else {
		result = fmt.Sprintf(`%s://%s%s`, u.Scheme, u.Host, u.EscapedPath())
	}
	return result
}
