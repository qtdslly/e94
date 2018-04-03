package middleware

import (
	"background/common/constant"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func GetDbPrepareHandler(dbName, dbSource string, enableLog bool) gin.HandlerFunc {
	db, err := gorm.Open(dbName, dbSource)
	if err != nil {
		return nil
	}

	if enableLog {
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}

	db.DB().SetMaxIdleConns(50)
	db.DB().SetMaxOpenConns(1500)

	return func(c *gin.Context) {
		c.Set(constant.ContextDb, db)
		c.Next()
	}
}
