package service

import (
	"fmt"
	"sort"
	"crypto/sha1"
	"encoding/hex"
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
