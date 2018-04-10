package model

import (
	"background/common/logger"
	"github.com/jinzhu/gorm"
)

func InitModel(db *gorm.DB) error {
	var err error

	err = initRealTimeStock(db)
	if err != nil {
		logger.Fatal("Init db real_time_stock failed, ", err)
		return err
	}

	err = initHoldStockInfo(db)
	if err != nil {
		logger.Fatal("Init db hold_stock_info failed, ", err)
		return err
	}

	err = initStockHistoryDataQ(db)
	if err != nil {
		logger.Fatal("Init db stock_history_data_q failed, ", err)
		return err
	}

	err = initStockList(db)
	if err != nil {
		logger.Fatal("Init db stock_list failed, ", err)
		return err
	}

	err = initTransPrompt(db)
	if err != nil {
		logger.Fatal("Init db trans_prompt failed, ", err)
		return err
	}

	err = initTransStockInfo(db)
	if err != nil {
		logger.Fatal("Init db trans_stock_info failed, ", err)
		return err
	}

	err = initDeepFallStock(db)
	if err != nil {
		logger.Fatal("Init db deep_fall_stock failed, ", err)
		return err
	}

	err = initSimulation(db)
	if err != nil {
		logger.Fatal("Init db simulation failed, ", err)
		return err
	}

	err = initTonghuashunSuggestion(db)
	if err != nil {
		logger.Fatal("Init db tonghuashun_suggestion failed, ", err)
		return err
	}

	err = initTonghuashunMainForceControl(db)
	if err != nil {
		logger.Fatal("Init db tonghuashun_main_force_control failed, ", err)
		return err
	}

	return err
}

// Do not call this method!!!!
func rebuildModel(db *gorm.DB) {
	dropRealTimeStock(db)
	dropHoldStockInfo(db)
	dropStockHistoryDataQ(db)
	dropStockList(db)
	dropTransPrompt(db)
	dropTransStockInfo(db)
	dropDeepFallStock(db)
	dropSimulation(db)
	dropTonghuashunSuggestion(db)
	dropTonghuashunMainForceControl(db)
	InitModel(db)
}
