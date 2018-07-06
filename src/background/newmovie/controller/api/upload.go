package api

import (
	"net/http"
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"os"
	"io"
	"time"
	"path/filepath"
	"background/stock/config"
)

func Fileupload(c *gin.Context){
	//得到上传的文件
	file, header, err := c.Request.FormFile("image") //image这个是uplaodify参数定义中的   'fileObjName':'image'
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	//文件的名称
	filename := header.Filename

	relDir := time.Now().Format("/2006/01/02/15/04")
	relPath := filepath.Join(relDir, filename)

	if err := os.MkdirAll(filepath.Join(config.GetStorageRoot(), relDir), 0755); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	tmpPath := filepath.Join(config.GetStorageRoot(),relPath)
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_ ,err = io.Copy(tmpFile, file)
	tmpFile.Close()
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success,"file_name":relPath})
}