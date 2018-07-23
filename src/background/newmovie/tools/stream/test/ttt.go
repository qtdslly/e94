package main

import (
	"background/common/aes1"
	"background/common/logger"
	"fmt"
)
func main(){
	url := "http://223.110.243.142/PLTV/2510088/224/3221227171/1.m3u8"
	data,err := aes1.Encrypt([]byte(url))
	if err != nil{
		logger.Error(err)
		return
	}

	fmt.Println(data)
}
