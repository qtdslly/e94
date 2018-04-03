package fms

import (
	"common/constant"
	"common/logger"
	"common/uuid"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	fileKey      string
	fmsRootPath  string
	fmsTmpPath   string
	fmsVideoPath string
	fmsUrlPrefix string
)

// 分片文件数组
type filePathSlice []string

// 此数组内都是文件名为数字的分片文件 按文件名数字大小进行排序
func (c filePathSlice) Len() int { return len(c) }
func (c filePathSlice) Less(i, j int) bool {
	num1, _ := strconv.Atoi(strings.TrimSuffix(filepath.Base(c[i]), filepath.Ext(c[i])))
	num2, _ := strconv.Atoi(strings.TrimSuffix(filepath.Base(c[j]), filepath.Ext(c[j])))
	return num1 < num2
}
func (c filePathSlice) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

func InitFms(storageRoot, urlPrefix string, key string) {
	fmsRootPath = storageRoot
	fmsTmpPath = filepath.Join(fmsRootPath, constant.TmpStorage)
	fmsVideoPath = filepath.Join(fmsRootPath, constant.VideoTmpStorage)
	fmsUrlPrefix = urlPrefix
	fileKey = key
}

func FormatDateDir(t *time.Time) string {
	return t.Format("/2006/01/02")
}

func UploadFile(filename string, f multipart.File) (relPath string, err error) {
	// make sure path exists
	os.MkdirAll(fmsTmpPath, 0755)

	tmpPath := filepath.Join(fmsTmpPath, uuid.NewUUID().String()+filepath.Ext(filename))
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, f)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	relaPath, _ := filepath.Rel(fmsRootPath, tmpPath)
	relaPath = filepath.Clean("/" + relaPath)

	return relaPath, nil
}

func DownloadFile(url string) (relPath string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("bad response status [%s] !", resp.Status))
		logger.Error(err)
		return "", err
	}

	// close body read before return
	defer resp.Body.Close()

	// make sure path exists
	os.MkdirAll(fmsTmpPath, 0755)

	tmpPath := filepath.Join(fmsTmpPath, uuid.NewUUID().String()+filepath.Ext(url))
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	relaPath, _ := filepath.Rel(fmsRootPath, tmpPath)
	relaPath = filepath.Clean("/" + relaPath)

	return relaPath, nil
}

func UpdateFile(src, targetPath string) (newPath string, err error) {
	var relPath string

	// check if targetpath is relative path, return if not
	if !filepath.HasPrefix(targetPath, "/") && !filepath.HasPrefix(targetPath, "\\") {
		err := errors.New(fmt.Sprintf("invalid target path [%s], must starts with /", targetPath))
		logger.Error(err)
		return "", err
	}

	// check if targetpath is in tmp, return if yes
	if filepath.HasPrefix(targetPath, "/tmp") || filepath.HasPrefix(targetPath, "\\tmp") {
		err := errors.New(fmt.Sprintf("invalid target path [%s], should not in /tmp", targetPath))
		logger.Error(err)
		return "", err
	}

	if filepath.HasPrefix(src, "http") {
		// shouldn't come here
		relPath, err = DownloadFile(src)
		if err != nil {
			logger.Error(err)
			return "", err
		}
	} else {
		relPath = filepath.Clean("/" + src)
	}

	// the src is in /tmp folder, target is not, move the file from tmp to target
	absTargetPath := filepath.Join(fmsRootPath, targetPath)
	// make directory if not exists
	os.MkdirAll(filepath.Dir(absTargetPath), 0755)
	err = os.Rename(filepath.Join(fmsRootPath, relPath), absTargetPath)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	return targetPath, nil
}

func BuildTemporaryFile(ext string) (relPath string, err error) {
	// make sure path exists
	os.MkdirAll(fmsTmpPath, 0755)

	tmpPath := filepath.Join(fmsTmpPath, uuid.NewUUID().String()+ext)

	path, _ := filepath.Rel(fmsRootPath, tmpPath)
	path = filepath.Clean("/" + path)
	return path, nil
}

func GetAbsPath(relPath string) string {
	return filepath.Join(fmsRootPath, relPath)
}

func FileServeHandler(c *gin.Context) {
	path := c.Request.URL.Path
	absPath := filepath.Join(fmsRootPath, path)
	info, err := os.Stat(absPath)
	if err == nil && info.Mode().IsRegular() {
		http.ServeFile(c.Writer, c.Request, absPath)
		return
	}

	if filepath.HasPrefix(path, fmsUrlPrefix) {
		relPath, err := filepath.Rel(fmsUrlPrefix, path)
		if err == nil {
			absPath := filepath.Join(fmsRootPath, relPath)
			info, err := os.Stat(absPath)
			if err == nil && info.Mode().IsRegular() {
				http.ServeFile(c.Writer, c.Request, absPath)
				return
			}
		}
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func CheckChunk(chunk int, chunkSize int64, adminId uint32, fileMd5 string) bool {
	userTmpPath := filepath.Join(filepath.Join(fmsRootPath, "admins"), fmt.Sprintf(`%d`, adminId))
	if fileMd5 != "" {
		userTmpPath = filepath.Join(userTmpPath, fileMd5)
	}
	tmpPath := filepath.Join(userTmpPath, fmt.Sprintf(`%d`, chunk))

	f, err := os.Stat(tmpPath)
	if err != nil {
		return false
	}

	fileSize := f.Size()
	if fileSize == chunkSize {
		return true
	}
	return false
}

func MergeUserChunkFiles(filename string, adminId uint32, fileMd5 string) bool {
	files, err := GetUserChunkFiles(adminId, fileMd5)
	if err != nil {
		logger.Error(err)
		return false
	}
	userTmpPath := filepath.Join(filepath.Join(fmsRootPath, "admins"), fmt.Sprintf(`%d`, adminId))
	os.MkdirAll(userTmpPath, 0755)

	tmpPath := filepath.Join(userTmpPath, filename)
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	defer tmpFile.Close()
	for _, file := range files {
		chunkFile, err := os.OpenFile(filepath.Join(fmsRootPath, file), os.O_RDONLY, 0755)
		if err != nil {
			logger.Error(err)
			return false
		}
		//bReader := bufio.NewReader(chunkFile)
		defer chunkFile.Close()

		_, err = io.Copy(tmpFile, chunkFile)
		if err != nil {
			logger.Error(err)
			return false
		}
	}
	return true
}

func GetUserChunkFiles(adminId uint32, fileMd5 string) (list []string, err error) {
	userTmpPath := filepath.Join(filepath.Join(fmsRootPath, "admins"), fmt.Sprintf(`%d`, adminId), fileMd5)
	os.MkdirAll(userTmpPath, 0755)
	dirList, err := ioutil.ReadDir(userTmpPath)
	if err != nil {
		logger.Error(err)
		return list, err
	}

	for _, item := range dirList {
		if !item.IsDir() {
			tmpPath := filepath.Join(userTmpPath, item.Name())
			realPath, _ := filepath.Rel(fmsRootPath, tmpPath)
			list = append(list, realPath)
		}
	}
	sort.Sort(filePathSlice(list))

	return list, err
}

/*
	upload user tmp file
*/
func UploadUserTmpFile(filename string, f multipart.File, adminId uint32, path string) (relPath string, err error) {
	// make sure path exists
	userTmpPath := filepath.Join(filepath.Join(fmsRootPath, "admins"), fmt.Sprintf(`%d`, adminId))
	if path != "" {
		userTmpPath = filepath.Join(userTmpPath, path)
	}
	os.MkdirAll(userTmpPath, 0755)

	tmpPath := filepath.Join(userTmpPath, filename)

	// 如果文件存在，则覆盖
	var tmpFile *os.File
	tmpFile, err = os.OpenFile(tmpPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	defer func() {
		tmpFile.Close()
		// 如果文件上传失败 删除创建的临时文件
		if err != nil {
			os.Remove(tmpPath)
		}
	}()

	_, err = io.Copy(tmpFile, f)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	relaPath, _ := filepath.Rel(fmsRootPath, tmpPath)
	relPath = filepath.Clean("/" + relaPath)
	return relPath, nil
}

/*
	get user tmp files
*/
func GetUserTmpFiles(adminId uint32) (list []string, err error) {
	userTmpPath := filepath.Join(filepath.Join(fmsRootPath, "admins"), fmt.Sprintf(`%d`, adminId))
	os.MkdirAll(userTmpPath, 0755)
	dirList, err := ioutil.ReadDir(userTmpPath)
	if err != nil {
		logger.Error(err)
		return list, err
	}

	for _, item := range dirList {
		if !item.IsDir() {
			tmpPath := filepath.Join(userTmpPath, item.Name())
			relaPath, _ := filepath.Rel(fmsRootPath, tmpPath)
			list = append(list, "/"+relaPath)
		}
	}
	return list, err
}

/*
	delete user tmp files
*/
func DeleteUserTmpFiles(adminId uint32, fileNames []string) {
	userTmpPath := filepath.Join(filepath.Join(fmsRootPath, "admins"), fmt.Sprintf(`%d`, adminId))
	for _, filename := range fileNames {
		filePath := filepath.Join(userTmpPath, filename)
		os.Remove(filePath)
	}
}
