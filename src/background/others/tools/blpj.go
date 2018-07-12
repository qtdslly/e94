package main

import (
	"net/http"
	//"encoding/json"
	"background/common/logger"
	"io/ioutil"
	"background/wechart/config"
	//"bytes"
	"crypto/md5"
	"fmt"

	"golang.org/x/text/encoding/simplifiedchinese"

	"strings"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	//type UserInfo struct {
	//	UserName   string    `form:"username"`
	//	Pwd        string    `form:"password"`
	//}
	//
	//var userInfo UserInfo
	//userInfo.UserName = "admin"
	//userInfo.Pwd = "admin"
	//
	//bytesData, err := json.Marshal(userInfo)
	//if err != nil {
	//	logger.Error(err)
	//	return
	//}
	//logger.Debug(string(bytesData))
	//
	//reader := bytes.NewReader(bytesData)

	apiUrl := "http://great-cause.com.cn/portal/userlogin"

	requ, err := http.NewRequest("POST", apiUrl, strings.NewReader("username=bjzy&password=123456"))

	requ.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	requ.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3278.0 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(requ)
	if err != nil {
		logger.Error(err)
		return
	}

	defer resp.Body.Close()

	recv, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debug(string(recv))
}


func GetMd5(str string)(string){
	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str1
}


func DecodeToGBK1(text string) (string, error) {

	dst := make([]byte, len(text)*2)
	tr := simplifiedchinese.GB18030.NewDecoder()
	nDst, _, err := tr.Transform(dst, []byte(text), true)
	if err != nil {
		return text, err
	}

	return string(dst[:nDst]), nil
}