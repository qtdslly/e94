package util

import (
	"golang.org/x/text/encoding/simplifiedchinese"
)
func OnlyNumber(str string)(string){
	//只保留数字
	//48-57
	rTitle := ([]rune)(str)
	result := ""
	for _, m := range rTitle {
		if m >= 48 && m <= 57{
			result += string(m)
		}
	}
	return result
}


func DecodeToGBK(text string) (string, error) {

	dst := make([]byte, len(text)*2)
	tr := simplifiedchinese.GB18030.NewDecoder()
	nDst, _, err := tr.Transform(dst, []byte(text), true)
	if err != nil {
		return text, err
	}

	return string(dst[:nDst]), nil
}