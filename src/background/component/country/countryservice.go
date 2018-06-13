package country

import (
	"common/constant"
	"common/logger"
	"common/uuid"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
)

/*
	This method allows to initialize the area database from the data text.
*/
func InitCountryData(areaDataLink, storageRootPath string, db *gorm.DB) {
	// 检查表内是否有数据
	var count uint32
	if err := db.Raw("select 1 from country limit 1").Count(&count).Error; err == nil {
		return
	}

	// 若表内无数据 则先下载文件到目录，再加载数据
	absPath := filepath.Join(storageRootPath, constant.TmpStorage)
	os.MkdirAll(absPath, 0755)
	absFilePath := filepath.Join(absPath, uuid.NewUUID().String()+filepath.Ext(areaDataLink))
	if err := downloadTmpFile(areaDataLink, absFilePath); err != nil {
		logger.Error(err)
		return
	}

	loadCountryData(absFilePath, db)
}

// 下载临时文件
func downloadTmpFile(url, absPath string) error {
	// remove file first
	os.Remove(absPath)

	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("bad response status [%s] !", resp.Status))
		logger.Error(err)
		return err
	}

	// close body read before return
	defer resp.Body.Close()

	tmpFile, err := os.Create(absPath)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// 加载country.data数据到数据库
func loadCountryData(path string, db *gorm.DB) {
	var f *os.File
	var err error
	f, err = os.Open(path)
	if err != nil {
		logger.Error(err)
		return
	}

	tx := db.Begin()
	defer func() {
		f.Close()
		os.Remove(path) // 删除文件
		if err != nil {
			tx.Rollback()
		}
	}()

	var data []byte
	data, err = ioutil.ReadAll(f)
	if err != nil {
		logger.Error(err)
		return
	}

	var countries []*Country

	// 数据初始化
	err = json.Unmarshal(data, &countries)
	if err != nil {
		logger.Error(err)
		return
	}

	for _, country := range countries {
		if err = tx.Create(&country).Error; err != nil {
			logger.Error(err)
			return
		}
	}

	tx.Commit()
}
