package area

import (
	"bufio"
	"bytes"
	"common/constant"
	"common/logger"
	"common/uuid"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jinzhu/gorm"
)

/*
	This method allows to initialize the area database from the data text.
*/
func InitAreaData(areaDataLink, storageRootPath string, db *gorm.DB) {
	// 检查表内是否有数据
	var count uint32
	if err := db.Raw("select 1 from area limit 1").Count(&count).Error; err == nil {
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

	loadAreaData(absFilePath, db)
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

// 加载area.data数据到数据库
func loadAreaData(path string, db *gorm.DB) {
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

	// 数据初始化
	buf := bytes.NewReader(data)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)

		if len(words) != 2 {
			logger.Error(fmt.Sprintf("line [%s] doesn't have 2 words", line))
			continue
		}
		// and we can iterate on top of that
		area := Area{}
		area.Code = strings.TrimSpace(words[0])
		area.Name = strings.TrimSpace(words[1])

		if strings.HasSuffix(area.Code, "0000") {
		} else if strings.HasSuffix(area.Code, "00") {
			area.ParentCode = area.Code[:2] + "0000"
		} else {
			area.ParentCode = area.Code[:4] + "00"
		}

		// check if the area_code already exists, if yes then skip
		if err = tx.Where("code = ?", area.Code).First(&area).Error; err == nil {
			// if area_code found, then skip
			continue
		}

		if err = tx.Create(&area).Error; err != nil {
			tx.Rollback()
			logger.Error(err)
			return
		}

	}

	if err = scanner.Err(); err != nil {
		tx.Rollback()
		logger.Error(err)
		return
	}

	tx.Commit()
}

/*
	This method allows to retrive the sub areas for a given areas and the level limit
*/
func GetSubAreaCodes(code string, db *gorm.DB) ([]string, error) {
	var area Area
	area.Code = code
	if err := db.First(&area).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	subAreas, err := GetSubAreas(&area, 5, db)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var codes []string
	for _, subArea := range subAreas {
		codes = append(codes, subArea.Code)
	}
	return codes, nil
}

/*
	This method allows to retrive the sub areas for a given areas and the level limit
*/
func GetSubAreas(dbArea *Area, level int, db *gorm.DB) ([]*Area, error) {
	// the total areas including current area and its sub areas
	var dbSubAreas []*Area

	if level == 0 {
		dbSubAreas = append(dbSubAreas, dbArea)
		return dbSubAreas, nil
	}

	// the sub companies level 1 for current Company
	var dbSubAreasL1 []*Area
	if err := db.Where("parent_code = ?", dbArea.Code).Find(&dbSubAreasL1).Error; err != nil {
		return nil, err
	}

	for _, dbSubArea := range dbSubAreasL1 {
		currentDbSubAreas, err := GetSubAreas(dbSubArea, level-1, db)
		if err != nil {
			return nil, err
		}
		dbSubAreas = append(dbSubAreas, currentDbSubAreas...)
	}
	dbSubAreas = append(dbSubAreas, dbArea)
	return dbSubAreas, nil
}

// retrive the parent area codes for a given area_code.
func GetParentArea(areaCode string) string {
	areaCode = strings.TrimSpace(areaCode)
	if len(areaCode) != 6 {
		logger.Error(fmt.Sprintf("invalid area_code [%s] length", areaCode))
		return ""
	}

	if strings.HasSuffix(areaCode, "0000") {
		return ""
	} else if strings.HasSuffix(areaCode, "00") {
		return areaCode[:2] + "0000"
	} else {
		return areaCode[:4] + "00"
	}
}
