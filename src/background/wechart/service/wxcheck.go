package service

import (
	"fmt"
	"sort"
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"github.com/tidwall/gjson"
	"background/common/logger"
	"io/ioutil"
)

func Check(timeStamp int64,nonce,signature string)bool{
	token := "5a61efdc52411a670b9f7c9db0a5275b"
	tmpStr := []string{fmt.Sprint(timeStamp),nonce,token}
	sort.Strings(tmpStr)
	sortStr := tmpStr[0] + tmpStr[1] + tmpStr[2]

	r := sha1.Sum([]byte(sortStr))
	result := hex.EncodeToString(r[:])
	if result == signature{
		return true
	}
	return false
}


func AccessToken()string{
	apiUrl := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=wx508e9e50a737c414&secret=58ab0927aeea25867c88d61530b985e7"

	requ, err := http.NewRequest("GET",apiUrl,nil)

	resp, err := http.DefaultClient.Do(requ)
	if err != nil {
		logger.Debug("Proxy failed!")
		return ""
	}

	recv,err := ioutil.ReadAll(resp.Body)
	if err != nil{
		logger.Error(err)
		return ""
	}

	logger.Debug(string(recv))
	accessToken := gjson.Get(string(recv), "access_token")

	if !accessToken.Exists(){
		errCode := gjson.Get(string(recv), "errcode")
		errMsg := gjson.Get(string(recv), "errmsg")
		logger.Debug(errCode)
		logger.Debug(errMsg)
		return ""
	}
	return accessToken.String()
}