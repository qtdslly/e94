package cms

import (
	"net/http"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"os"
	"io"
	"background/stock/config"
)

func FileUpload(c *gin.Context){
	//得到上传的文件
	file, header, err := c.Request.FormFile("file") //image这个是uplaodify参数定义中的   'fileObjName':'image'
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	//文件的名称
	filename := header.Filename

	tmpFile, err := os.Create(config.GetStorageRoot() + filename)
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
	c.String(http.StatusOK,"文件上传成功")
}