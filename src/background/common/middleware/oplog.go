package middleware

import (
	"background/common/constant"
	"background/common/logger"
	"background/common/util"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

/*
	Warning: keep this consistent with ams/model's OperationLog
*/
type OperationLog struct {
	Id            uint32    `gorm:"primary_key" json:"id"`
	Table         string    `gorm:"size:64" json:"table"`
	Type          uint32    `json:"type"`
	Method        string    `gorm:"size:16" json:"method"`
	RequestUrl    string    `gorm:"size:255" json:"request_url"`
	RequestBody   string    `gorm:"type:longtext" json:"request_body"`
	OldValue      string    `gorm:"type:longtext" json:"old_value"`
	NewValue      string    `gorm:"type:longtext" json:"new_value"`
	Error         string    `gorm:"size:255" json:"error"`
	Operator      string    `gorm:"size:32" json:"operator"`
	OperatorId    uint32    `json:"operator_id"`
	OperatorIp    string    `gorm:"size:64" json:"operator_ip"`
	ExecutionTime uint32    `json:"execution_time"` // millisecond
	Status        uint32    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

/*
	This middleware allows to generate a operation log for some important api operation.
*/
func OperationLogHandler(amsBindAddr string, token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		var elapsed uint32 // millisecond
		elapsed = uint32(time.Since(start) / time.Millisecond)

		var table string
		var typ uint32
		var err string

		tab, exists := c.Get(constant.ContextTableName)
		if exists && tab != nil {
			table = tab.(string)
		} else {
			// if table name is not set, then return
			return
		}

		opType, exists := c.Get(constant.ContextOperationType)
		if exists && tab != nil {
			typ = uint32(opType.(int))
		}

		e, exists := c.Get(constant.ContextError)
		if exists && e != nil {
			err = e.(string)
		}

		var adminId uint32
		ad, exists := c.Get(constant.ContextAdminId)
		if exists && ad != nil {
			adminId = ad.(uint32)
		}

		var adminName string
		adName, exists := c.Get(constant.ContextAdminName)
		if exists && adName != nil {
			adminName = adName.(string)
		}

		loginIp := c.ClientIP()
		requestUrl := c.Request.URL.RequestURI()

		requestBody, _ := c.Get(constant.ContextRequestBody)
		oldValue, _ := c.Get(constant.ContextOldValue)
		newValue, _ := c.Get(constant.ContextNewValue)
		status := uint32(c.Writer.Status())

		// collect all useful information and create a operation_log.
		var opLog OperationLog
		opLog.Table = table
		opLog.Method = c.Request.Method
		opLog.RequestUrl = requestUrl
		opLog.Type = typ

		if requestBody != nil {
			b, err := json.Marshal(requestBody)
			if err != nil {
				logger.Error(err)
				return
			}
			opLog.RequestBody = string(b)
		}

		if oldValue != nil {
			b, err := json.Marshal(oldValue)
			if err != nil {
				logger.Error(err)
				return
			}
			opLog.OldValue = string(b)
		}

		if newValue != nil {
			b, err := json.Marshal(newValue)
			if err != nil {
				logger.Error(err)
				return
			}
			opLog.NewValue = string(b)
		}

		opLog.Error = err
		opLog.Operator = adminName
		opLog.OperatorId = adminId
		opLog.OperatorIp = loginIp
		opLog.ExecutionTime = elapsed
		opLog.Status = status

		timestamp := time.Now().Unix()
		params := make(map[string]interface{})
		opJson, _ := json.Marshal(opLog)
		json.Unmarshal(opJson, &params)
		params["timestamp"] = timestamp

		moduleSign := MakeModuleSignature(params, constant.ModuleSalt, timestamp)

		//logger.Debug("ModuleSign:", moduleSign)

		type response struct {
			ErrCode int `json:"err_code"`
		}
		var resp response
		api := amsBindAddr
		api += "/ams/operationlog/create"
		if err := util.Post(api, &params, &resp, map[string]string{constant.AccessKey: token, constant.ModuleSignature: moduleSign}); err != nil {
			logger.Error(err)
		}
		return
	}
}
