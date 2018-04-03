package util

import (
	"background/common/logger"

	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

/*
	this method allows to send a request and bind the response to the provided json model.
*/
func Get(url string, model interface{}, headers map[string]string) error {
	req, err := http.NewRequest("GET", url, nil)
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("bad response status [%s] !", resp.Status))
		logger.Error(err)
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return err
	}
	// CAN BE REMOVED
	// logger.Debug(string(b))
	if err = json.Unmarshal(b, model); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

/*
	this method allows to send a request and bind the response to the provided xml model.
*/
func GetXml(url string, model interface{}, headers map[string]string) error {
	req, err := http.NewRequest("GET", url, nil)
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("bad response status [%s] !", resp.Status))
		logger.Error(err)
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return err
	}
	// CAN BE REMOVED
	// logger.Debug(string(b))
	if err = xml.Unmarshal(b, model); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

/*
	This method allows to post a json input and retrieve the response.
*/
func Post(url string, input interface{}, output interface{}, headers map[string]string) error {
	b, err := json.Marshal(input)
	if err != nil {
		logger.Error(err)
		return err
	}

	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json; charset=UTF-8")

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("bad response status [%s] !", resp.Status))
		logger.Error(err)
		return err
	}

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return err
	}

	if err = json.Unmarshal(out, output); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

/*
	This method allows to post a form input and retrieve the response.
*/
func PostFile(url string, formFile string, filename string, input []byte, output interface{}, headers map[string]string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile(formFile, filename)
	if err != nil {
		logger.Error("error writing to buffer", err)
		return err
	}
	fileWriter.Write(input)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		logger.Error(err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("bad response status [%s] !", resp.Status))
		logger.Error(err)
		return err
	}

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return err
	}

	if err = json.Unmarshal(out, output); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
