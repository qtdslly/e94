package model

import (
	"background/common/logger"
	"github.com/jinzhu/gorm"
)

func InitModel(db *gorm.DB) error {
	var err error

	err = initMovie(db)
	if err != nil {
		logger.Fatal("Init db movie failed, ", err)
		return err
	}
	err = initTopSearch(db)
	if err != nil {
		logger.Fatal("Init db top_search failed, ", err)
		return err
	}

	err = initTop(db)
	if err != nil {
		logger.Fatal("Init db top failed, ", err)
		return err
	}

	err = initAdmin(db)
	if err != nil {
		logger.Fatal("Init db admin failed, ", err)
		return err
	}

	err = initInstallation(db)
	if err != nil {
		logger.Fatal("Init db installation failed, ", err)
		return err
	}

	err = initKvStore(db)
	if err != nil {
		logger.Fatal("Init db kv_store failed, ", err)
		return err
	}

	err = initVideo(db)
	if err != nil {
		logger.Fatal("Init db video failed, ", err)
		return err
	}

	err = initEpisode(db)
	if err != nil {
		logger.Fatal("Init db episode failed, ", err)
		return err
	}

	err = initPlayUrl(db)
	if err != nil {
		logger.Fatal("Init db play_url failed, ", err)
		return err
	}

	err = initRecommend(db)
	if err != nil {
		logger.Fatal("Init db recommend failed, ", err)
		return err
	}
	err = initResourceGroup(db)
	if err != nil {
		logger.Fatal("Init db resource_group failed, ", err)
		return err
	}

	err = initNotification(db)
	if err != nil {
		logger.Fatal("Init db notification failed, ", err)
		return err
	}

	err = initStream(db)
	if err != nil {
		logger.Fatal("Init db stream failed, ", err)
		return err
	}

	err = initUser(db)
	if err != nil {
		logger.Fatal("Init db user failed, ", err)
		return err
	}

	err = initContentAction(db)
	if err != nil {
		logger.Fatal("Init db content_action failed, ", err)
		return err
	}

	err = initUserStream(db)
	if err != nil {
		logger.Fatal("Init db user_stream failed, ", err)
		return err
	}

	err = initUserOpinion(db)
	if err != nil {
		logger.Fatal("Init db user_opinion failed, ", err)
		return err
	}

	err = initUserWant(db)
	if err != nil {
		logger.Fatal("Init db user_want failed, ", err)
		return err
	}

	err = initFile(db)
	if err != nil {
		logger.Fatal("Init db file failed, ", err)
		return err
	}

	err = initApp(db)
	if err != nil {
		logger.Fatal("Init db app failed, ", err)
		return err
	}

	err = initVersion(db)
	if err != nil {
		logger.Fatal("Init db version failed, ", err)
		return err
	}

	err = initUpgrade(db)
	if err != nil {
		logger.Fatal("Init db version failed, ", err)
		return err
	}

	err = initTag(db)
	if err != nil {
		logger.Fatal("Init db tag failed, ", err)
		return err
	}


	err = initActivity(db)
	if err != nil {
		logger.Fatal("Init db activity failed, ", err)
		return err
	}
	return err
}

// Do not call this method!!!!
func rebuildModel(db *gorm.DB) {
	dropMovie(db)
	dropTopSearch(db)
	dropAdmin(db)
	dropInstallation(db)
	dropVideo(db)
	dropEpisode(db)
	dropPlayUrl(db)
	dropRecommend(db)
	dropResourceGroup(db)
	dropNotification(db)
	dropUser(db)
	dropContentAction(db)
	dropUserStream(db)
	dropUserOpinion(db)
	dropFile(db)
	dropUserWant(db)
	dropApp(db)
	dropVersion(db)
	dropTag(db)
	dropActivity(db)
	dropTop(db)

	InitModel(db)
}
