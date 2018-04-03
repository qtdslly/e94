package util

import (
	"background/common/constant"
	"background/common/logger"

	"errors"
	"fmt"
	"strings"
	"strconv"
)

// this method allows to call resource server and upload a file
func Upload(url, resBindAddr string) (string, error) {
	type request struct {
		Url string `json:"url"`
	}

	type response struct {
		ErrCode int    `json:"err_code"`
		Data    string `json:"data"`
	}

	url = strings.TrimSpace(url)
	if url == "" {
		return "", errors.New("Invalid upload url.")
	}

	apiUrl := resBindAddr + "/res/file/save"
	var resp response
	if err := Post(apiUrl, &request{Url: url}, &resp, nil); err != nil {
		logger.Error(err, url)
		return "", err
	}

	if resp.ErrCode != constant.Success {
		logger.Error("fail to save: ", url)
		err := fmt.Errorf("err_code [%d]", resp.ErrCode)
		return "", err
	}

	return resp.Data, nil
}

func UploadFileBytes(filename string, data []byte, resBindAddr string) (string, error) {
	type response struct {
		ErrCode int    `json:"err_code"`
		Data    string `json:"data"`
	}
	logger.Printf(filename)
	logger.Printf(resBindAddr)
	logger.Printf(strconv.Itoa(len(data)))
	apiUrl := resBindAddr + "/res/file/upload"
	var resp response
	if err := PostFile(apiUrl, "file", filename, data, &resp, nil); err != nil {
		logger.Error(err, apiUrl)
		return "", err
	}

	if resp.ErrCode != constant.Success {
		logger.Error("fail to save: ", apiUrl)
		err := fmt.Errorf("err_code [%d]", resp.ErrCode)
		return "", err
	}

	return resp.Data, nil
}
