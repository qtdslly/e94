package util

import (
	"background/common/constant"
	"regexp"
)

/*
	This method allows to extract device model.
*/
func GetDeviceModel(osType uint32, ua string) string {
	if osType == constant.OsTypeAndroid {
		return extract(ua, `;\s?([^;]+?)\s?(Build)?/`)
	} else if osType == constant.OsTypeIos {
		return extract(ua, `\((.*?);`)
	}
	return ""
}

/*
	This method allows to extract device model.
*/
func extract(ua string, pattern string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return ""
	}
	ret := re.FindSubmatch([]byte(ua))
	if len(ret) > 1 {
		return string(ret[1])
	}
	return ""
}
