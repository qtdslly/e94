package main

import (
	"net/http"
	"background/common/logger"
	"strings"
	"fmt"
	"path/filepath"
	"os"
	"io"
	"errors"
)

func main(){
	url := "http://pcvideogs.titan.mgtv.com/c1/2017/10/25_0/13DC5998B3400F05936BF120622D468F_20171025_1_1_1246.mp4?arange=0&scid=25012&pm=cRfgigekCgR9hEt8DSF68gd6d8EFscG7p3lwkZoWX91KGL2NQKD9zLsfHcImHMjVa7D1cKkGawGjlGjI4hOndMwUGS31FA~knEzBAy9SiWFQMU73WgpFgxZTbyLicHUwBxBkkRjd347marqzXi~cxLcn2xfYCrM0WgvgiOmwBt84N8e0aQgfWAtmBlI3MjSK0_PeXIq68VZspW7~3G2FW6~jcDSJCePNjb5XTEpis02B8MohyOHK2VnRua2WVHRSpj6jrKGYAVvpLk8CUZjYnZRmg~WfV3zTRd3sDolr_BoKiPbdWmqCE9727oY3OTQGrR5V2ubynOUYLiCrweR8n4FMr8_RKS2si8Z2bQcYmMJAomsdwBqN~DqAlwFqg7Oo3q~zNWQhmg~_ycjHsvYfOmk7j4c-&vcdn=0"
	fileName := "侏罗纪世界.mp4"

	DownloadFile(url,fileName)
}


func DownloadFile(requrl string, filename string) (int64, error) {
	//no timeout
	client := http.Client{}

	resp, err := client.Get(requrl)
	if err != nil {
	logger.Error(err)
	return 0, err
	}

	// close body read before return
	defer resp.Body.Close()

	// should not save html content as file
	if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
	err = errors.New("invalid response content type")
	logger.Error(err)
	return 0, err
	}

	if resp.StatusCode != http.StatusOK {
	err = errors.New(fmt.Sprintf("Fail to download: [%s]", requrl))
	logger.Error(err)
	return 0, err
	}

	tmpPath := filepath.Join("F:/ffmpeg/download", filename)
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
	logger.Error(err)
	return 0, err
	}

	bytes, err := io.Copy(tmpFile, resp.Body)
	tmpFile.Close()
	if err != nil {
	logger.Error(err)
	return 0, err
	}

	logger.Debug(tmpPath)
	return bytes, nil
}

