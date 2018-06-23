package main

import (
	"background/common/logger"
	"background/videodownload/config"
	"background/videodownload/model"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"flag"
	"strings"
	"errors"
	"net/http"
	"fmt"
	"path/filepath"
	"os"
	"io"
)

func main(){
	configPath := flag.String("conf", "../config/config.json", "Config file path")

	flag.Parse()

	logger.SetLevel(config.GetLoggerLevel())

	err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Error("Config Failed!!!!", err)
		return
	}

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}

	db.LogMode(true)
	model.InitModel(db)
	logger.SetLevel(config.GetLoggerLevel())

	charter := 1
	for {
		page := 1
		for {
			if page > 50{
				break
			}
			//http://mangafreak.com/images/manga/wrestle_the_under_ground/9/2.jpg
			DownloadFile("http://mangafreak.com/images/manga/wrestle_the_under_ground/" + fmt.Sprint(charter) + "/" + fmt.Sprint(page) + ".jpg", fmt.Sprintf("%03d",charter) + "-" + fmt.Sprintf("%02d",page) + ".jpg")
			page++
		}
		charter++
		if charter == 12{
			break
		}
	}


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

	tmpPath := filepath.Join("D:/data/cartoon/wrestle_the_under_ground/", filename)
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
