package cache

import (
	"cms/model"
	"cms/service"
	"common/constant"
	"common/logger"

	"github.com/jinzhu/gorm"
)

func getVideo(appId, id uint32, db *gorm.DB) (*model.Video, error) {
	var err error

	var video model.Video
	if err = db.Where("id = ?", id).First(&video).Error; err != nil {
		logger.Error(err)
		return nil, err

	}

	video.Cdns, err = GetDefaultCdn(appId, constant.MediaTypeVideo, video.Id, db)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &video, nil
}

// this method allows to fetch the video from cache, if not found then get it from db and update cache.
func GetVideo(appId, id uint32, db *gorm.DB) (*model.Video, error) {
	key := service.GetCacheKey("video", appId, 0, constant.MediaTypeVideo, id, "language_code", "zh")
	var video model.Video
	err := service.GetCacheObject(key, &video, func() (interface{}, error) {
		// if not found, then load from db
		v, err := getVideo(appId, id, db)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		return v, nil
	})

	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return &video, nil
}

func GetDefaultCdn(appId uint32, contentType uint8, contentId uint32, db *gorm.DB) ([]*model.Cdn, error) {
	key := service.GetCacheKey("cdn", appId, 0, uint32(contentType), contentId)
	var cdns []*model.Cdn
	err := service.GetCacheObject(key, &cdns, func() (interface{}, error) {
		cs, err := service.GetDefaultCdn(appId, contentType, contentId, db)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		return cs, nil
	})

	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return cdns, nil
}
