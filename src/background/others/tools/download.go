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
	url := "ftp://ygdy8:ygdy8@yg45.dydytt.net:8195/阳光电影www.ygdy8.com.我们诞生在中国.BD.720p.国英双语双字.mkv"
	fileName := "我们诞生在中国.mp4"

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

