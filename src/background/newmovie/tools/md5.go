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
		fmt.Println("%x", md5h.Sum([]byte(""))) //md5
	}else{
		fmt.Println(inerr)
	}
}