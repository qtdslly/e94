package main

import (
	"background/shortvideo/config"
	"background/common/logger"

	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {
	configPath := flag.String("conf", "../config/config.json", "Config file path")

	flag.Parse()

	logger.SetLevel(config.GetLoggerLevel())

	err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Error("Config Failed!!!!", err)
		return
	}

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}

	db.LogMode(true)
	logger.SetLevel(config.GetLoggerLevel())

	for {
		if !GetDouYin() {
			break
		}
	}
}


func GetDouYin() bool {

	apiurls := []string{
		"https://aweme-eagle-hl.snssdk.com/aweme/v1/feed/?version_code=5.8.0&pass-region=1&pass-route=1&js_sdk_version=1.13.3.0&app_name=aweme&vid=505ECF2B-061B-40AF-AA0E-2F4DB573FCAA&app_version=5.8.0&device_id=51632502948&channel=App%20Store&mcc_mnc=46001&aid=1128&screen_width=1242&openudid=0912a38b4938c6d19a0e2f0a745339deb70bbe82&os_api=18&ac=WIFI&os_version=12.2&device_platform=iphone&build_number=58010&device_type=iPhone8,2&iid=70375126620&idfa=F7835587-9CBD-4B92-8C6E-2813811E7B5F&volume=0.00&count=6&longitude=114.3404602033706&feed_style=0&filter_warn=0&address_book_access=1&user_id=102900047091&type=0&latitude=30.55565496650956&gps_access=3&pull_type=2&max_cursor=0&mas=01775e4a637fafd29a0e52c72f5763221a8856dae5bf755107eeb7&as=a2d5cc4cf7161cc9222904&ts=1556269415",
	}

	for _, apiurl := range apiurls {
		requ, err := http.NewRequest("GET", apiurl, nil)
		requ.Header.Add("authority", "aweme-eagle-hl.snssdk.com")
		requ.Header.Add("User-Agent", "Aweme 5.8.0 rv:58010 (iPhone; iOS 12.2; zh_CN) Cronet")

		requ.Header.Add("authority","aweme-eagle-hl.snssdk.com")
		requ.Header.Add("scheme", "https")
		requ.Header.Add("path", "/aweme/v1/feed/?version_code=5.8.0&pass-region=1&pass-route=1&js_sdk_version=1.13.3.0&app_name=aweme&vid=505ECF2B-061B-40AF-AA0E-2F4DB573FCAA&app_version=5.8.0&device_id=51632502948&channel=App%20Store&mcc_mnc=46001&aid=1128&screen_width=1242&openudid=0912a38b4938c6d19a0e2f0a745339deb70bbe82&os_api=18&ac=WIFI&os_version=12.2&device_platform=iphone&build_number=58010&device_type=iPhone8,2&iid=70375126620&idfa=F7835587-9CBD-4B92-8C6E-2813811E7B5F&volume=0.00&count=6&longitude=114.3404602033706&feed_style=0&filter_warn=0&address_book_access=1&user_id=102900047091&type=0&latitude=30.55565496650956&gps_access=3&pull_type=2&max_cursor=0&mas=011b451633a6db0d55492c9d02dbda679585e13cc8105fe47fd6a8&as=a2a58c1c91161c8a520118&ts=1556269665")
		requ.Header.Add("x-tt-token", "00576f8767f6b42902ce9e47ffe26c17828ef0bf17f3bd1df1541031f4f3c77e895f4ee87c2703ff3ba5bd83cc98ffc56e34")
		requ.Header.Add("sdk-version", "1")
		requ.Header.Add("user-agent", "Aweme 5.8.0 rv:58010 (iPhone; iOS 12.2; zh_CN) Cronet")
		requ.Header.Add("x-ss-tc", "0")
		requ.Header.Add("accept-encoding", "gzip, deflatev")
		requ.Header.Add("cookie", "odin_tt=9f32edf3c45ec68ecb9fc073cbccb2df7e1ab3c1ae4032efca3d80306c6b059db9d4b709901f208762e0146985b098fb49ebbfece7734586d6dc3fc5c9fbccb4")
		requ.Header.Add("cookie", "sid_guard=576f8767f6b42902ce9e47ffe26c1782%7C1555049213%7C5184000%7CTue%2C+11-Jun-2019+06%3A06%3A53+GMT")
		requ.Header.Add("cookie", "uid_tt=7904ea4d317494250139fb9055d0b181")
		requ.Header.Add("cookie", "sid_tt=576f8767f6b42902ce9e47ffe26c1782")
		requ.Header.Add("cookie", "sessionid=576f8767f6b42902ce9e47ffe26c1782")
		requ.Header.Add("cookie", "install_id=70375126620")
		requ.Header.Add("cookie", "ttreq=1$6702b323df50dfa2e1017a2aff8b59a909601dd7")
		requ.Header.Add("x-khronos", "1556269665")
		requ.Header.Add("x-gorgon", "8300000000004f3673dab28e07e3e8519edc85d18bc9efbaa7c3")


		resp, err := http.DefaultClient.Do(requ)
		if err != nil {
			logger.Debug("Proxy failed!")
			return false
		}

		recv, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error(err)
			return false
		}

		data,err := DecodeToGBK(string(recv))
		if err != nil{
			logger.Error(err)
			return false
		}
		logger.Debug(data)
		//logger.Debug(string(recv))




	}

	return true
}

func DownloadDouYinFile(requrl string, filename string) (int64, error) {
	//no timeout
	client := http.Client{}

	resp, err := client.Get(requrl)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	// close body read before return
	defer resp.Body.Close()

	// should not save html content as file
	if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		err = errors.New("invalid response content type")
		logger.Error(err)
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Fail to download: [%s]", requrl))
		logger.Error(err)
		return 0, err
	}

	tmpPath := filepath.Join(config.GetStorageRoot(), filename)
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	bytes, err := io.Copy(tmpFile, resp.Body)
	tmpFile.Close()
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	logger.Debug(tmpPath)
	return bytes, nil
}



func DecodeToGBK(text string) (string, error) {

	dst := make([]byte, len(text)*2)
	tr := simplifiedchinese.GB18030.NewDecoder()
	nDst, _, err := tr.Transform(dst, []byte(text), true)
	if err != nil {
		return text, err
	}

	return string(dst[:nDst]), nil
}