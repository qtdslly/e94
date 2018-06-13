package setting

import (
	"background/common/constant"
	"background/common/logger"
	"background/component/kv"
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type LibrarySetting struct {
	CategoryId uint32             `json:"category_id"`
	Style      *LibraryStyle      `json:"style"`
	Properties []*LibraryProperty `json:"properties"`
}

type LibraryStyle struct {
	NumPerRow          int  `json:"num_per_row"`
	DsplayScore        bool `json:"display_score"`
	DisplayTitle       bool `json:"display_title"`
	DisplayDescription bool `json:"display_description"`
}

type LibraryProperty struct {
	Id   uint32        `json:"id"`
	Name string        `json:"name"`
	Sort uint32        `json:"sort"`
	Tags []*LibraryTag `json:"tags"`
}

type LibraryTag struct {
	Id   uint32 `json:"id"`
	Name string `json:"name"`
	Sort uint32 `json:"sort"`
}

// 获取所有的设置
func GetLibrarySetting(appId, versionId uint32, db *gorm.DB) ([]*LibrarySetting, error) {
	value, err := kv.GetValueForKey(appId, versionId, constant.LibrarySettingKey, db)
	if err != nil {
		return nil, err
	}

	if value == "" {
		return []*LibrarySetting{}, nil
	}

	var librarySet []*LibrarySetting
	err = json.Unmarshal([]byte(value), &librarySet)
	if err != nil {
		return nil, err
	}
	return librarySet, nil
}

func SetLibrarySetting(appId, versionId uint32, librarySet []*LibrarySetting, db *gorm.DB) error {
	value, err := json.Marshal(librarySet)
	if err != nil {
		return err
	}

	err = kv.SetValueForKey(appId, versionId, constant.LibrarySettingKey, string(value), db)
	if err != nil {
		return err
	}

	return nil
}

// 根据category_id获取设置
func GetLibrarySettingByCid(appId, versionId, categoryId uint32, db *gorm.DB) *LibrarySetting {
	libraryCache, err := GetLibrarySetting(appId, versionId, db)
	if err != nil {
		logger.Error(err)
		return &LibrarySetting{CategoryId: categoryId}
	}

	for _, item := range libraryCache {
		if item.CategoryId == categoryId {
			return item
		}
	}

	return &LibrarySetting{CategoryId: categoryId}
}

// 更新单条记录配置
func SetLibrarySettingByItem(appId, versionId uint32, set *LibrarySetting, db *gorm.DB) error {
	libraryCache, err := GetLibrarySetting(appId, versionId, db)
	if err != nil {
		logger.Error(err)
		return err
	}

	exist := false
	for index, item := range libraryCache {
		if item.CategoryId == set.CategoryId {
			exist = true
			libraryCache[index] = set
			break
		}
	}
	if !exist {
		libraryCache = append(libraryCache, set)
	}

	err = SetLibrarySetting(appId, versionId, libraryCache, db)
	return err
}

// 根据app_id和category_id获取设置
func GetLibrarySettingByApp(appId, versionId uint32, db *gorm.DB) []*LibrarySetting {
	libraryCache, err := GetLibrarySetting(appId, versionId, db)
	if err != nil {
		logger.Error(err)
		return libraryCache
	}
	return libraryCache
}

// 删除配置
func DeleteLibrarySettingByVersion(versionId uint32, db *gorm.DB) error {
	return kv.DeleteValueForKeyByVersion(versionId, constant.LibrarySettingKey, db)
}

// 删除配置
func DeleteLibrarySettingByApp(appId uint32, db *gorm.DB) error {
	return kv.DeleteValueForKeyByApp(appId, constant.LibrarySettingKey, db)
}
