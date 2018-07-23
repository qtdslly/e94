package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)


func aaa(path string) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Open", err)
		return
	}

	defer f.Close()

	body, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("ReadAll", err)
		return
	}

	md5.Sum(body)
	fmt.Printf("%x\n", md5.Sum(body))
}

func bbb(path string) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Open", err)
		return
	}

	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		fmt.Println("Copy", err)
		return
	}

	md5hash.Sum(nil)
	//fmt.Printf("%x\n", md5hash.Sum(nil))
}

func main() {

}