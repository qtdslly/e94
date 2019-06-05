package api

import (
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"net/http"

	"time"
	"fmt"
	"background/wdkq/config"
	"os"

	"github.com/tealeg/xlsx"

	"strings"
	"strconv"
)

func SavePrintHandler(c *gin.Context) {

	logger.SetLevel(config.GetLoggerLevel())

	type param struct {
		AmtNo    string `form:"amt_no" json:"amt_no"`
		Name     string `form:"name" json:"name"`
		Sex      string `form:"sex" json:"sex"`
		Phone    string `form:"phone" json:"phone"`
		Birthday string `form:"birthday" json:"birthday"`
		Doctor   string `form:"doctor" json:"doctor"`
		Time     string `form:"time" json:"time"`
		GhAmt    float64 `form:"gh_amt" json:"gh_amt"`
		XyAmt    float64 `form:"xy_amt" json:"xy_amt"`
		FsAmt    float64 `form:"fs_amt" json:"fs_amt"`
		ZlAmt    float64 `form:"zl_amt" json:"zl_amt"`
		TotalAmt float64 `form:"total_amt" json:"total_amt"`
		SsAmt    float64 `form:"ss_amt" json:"ss_amt"`
		ZllAmt   float64 `form:"zll_amt" json:"zll_amt"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	var err error

	n := time.Now();

	fileName := fmt.Sprintf("%s%02d%02d%02d.xlsx", config.GetFilePath(), n.Year(), n.Month(), n.Day())

	if _, err := os.Stat(fileName); err != nil {
		if !os.IsExist(err) {
			file := xlsx.NewFile()

			sheet, err := file.AddSheet("预交金记录")
			if err != nil {
				logger.Error(err)
				c.JSON(http.StatusOK, gin.H{"err_code": constant.Failure, "err_msg": "打开sheet失败:" + err.Error()})
				return
			}

			row := sheet.AddRow()
			row.SetHeightCM(1)

			cell := row.AddCell()
			cell.Value = "共0笔"

			cell = row.AddCell()
			cell.Value = "合计0.0"

			row = sheet.AddRow()
			row.SetHeightCM(0.8)
			cell = row.AddCell()
			cell.Value = "收据号"

			cell = row.AddCell()
			cell.Value = "病人姓名"

			cell = row.AddCell()
			cell.Value = "性别"

			cell = row.AddCell()
			cell.Value = "电话"

			cell = row.AddCell()
			cell.Value = "生日"

			cell = row.AddCell()
			cell.Value = "接诊医生"

			cell = row.AddCell()
			cell.Value = "挂号费"

			cell = row.AddCell()
			cell.Value = "西药费"

			cell = row.AddCell()
			cell.Value = "放射费"

			cell = row.AddCell()
			cell.Value = "治疗费"

			cell = row.AddCell()
			cell.Value = "合计"

			cell = row.AddCell()
			cell.Value = "实收"

			cell = row.AddCell()
			cell.Value = "找零"

			cell = row.AddCell()
			cell.Value = "交易时间"

			if err := file.Save(fileName); err != nil {
				logger.Error(err);
				c.JSON(http.StatusOK, gin.H{"err_code": constant.Failure, "err_msg": "创建excel文件失败:" + err.Error()})
				return
			}
		}
	}

	file, err := xlsx.OpenFile(fileName)
	if err != nil {
		logger.Error(err);
		c.JSON(http.StatusOK, gin.H{"err_code": constant.Failure, "err_msg": "打开excel文件失败:" + err.Error()})
		return
	}

	sheet := file.Sheets[0]

	row1 := sheet.Rows[0]

	value := row1.Cells[0].Value
	sCount := strings.Replace(value, "共:", "", -1)
	sCount = strings.Replace(sCount, "笔", "", -1)
	count, _ := strconv.Atoi(sCount)
	row1.Cells[0].Value = fmt.Sprintf("共:%d笔", count + 1)

	value = row1.Cells[1].Value
	sCount = strings.Replace(value, "合计:", "", -1)

	amt, _ := strconv.ParseFloat(sCount, 64)
	row1.Cells[1].Value = fmt.Sprintf("合计:%.2f", amt + p.TotalAmt)

	row := sheet.AddRow()
	row.SetHeightCM(0.8)
	cell := row.AddCell()
	cell.Value = p.AmtNo

	cell = row.AddCell()
	cell.Value = p.Name

	cell = row.AddCell()
	cell.Value = p.Sex

	cell = row.AddCell()
	cell.Value = p.Phone

	cell = row.AddCell()
	cell.Value = p.Birthday

	cell = row.AddCell()
	cell.Value = p.Doctor

	cell = row.AddCell()
	cell.Value = fmt.Sprintf("%.2f", p.GhAmt)

	cell = row.AddCell()
	cell.Value = fmt.Sprintf("%.2f", p.XyAmt)

	cell = row.AddCell()
	cell.Value = fmt.Sprintf("%.2f", p.FsAmt)

	cell = row.AddCell()
	cell.Value = fmt.Sprintf("%.2f", p.ZlAmt)

	cell = row.AddCell()
	cell.Value = fmt.Sprintf("%.2f", p.TotalAmt)

	cell = row.AddCell()
	cell.Value = fmt.Sprintf("%.2f", p.SsAmt)

	cell = row.AddCell()
	cell.Value = fmt.Sprintf("%.2f", p.ZllAmt)

	cell = row.AddCell()
	cell.Value = p.Time

	err = file.Save(fileName)
	if err != nil {
		logger.Error(err);
		c.JSON(http.StatusOK, gin.H{"err_code": constant.Failure, "err_msg": "保存excel文件失败:" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}
