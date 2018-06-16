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
	url := "http://video.jiagew762.com:8091/20180515/B1WE78vo/index.m3u8"
	fileName := "a1.mp4"

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

