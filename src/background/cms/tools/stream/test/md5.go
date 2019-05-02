package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
)


func main() {
	testFile := "/root/data/storage/live.apk"
	log.Println(testFile)
	file, inerr := os.Open(testFile)
	if inerr == nil {
		md5h := md5.New()
		io.Copy(md5h, file)
		fmt.Printf("%x", md5h.Sum([]byte(""))) //md5
	}else{
		fmt.Println(inerr)
	}
}
//67 174 86 56 168 141 206 30 171 164 65 219 90 107 53 124
//241 201 139 185 162 92 212 54 166 225 26 106 69 38 59 135