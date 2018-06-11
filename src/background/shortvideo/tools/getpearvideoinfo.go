package main

import (
	"background/shortvideo/config"
	"background/shortvideo/model"
	"background/common/logger"

	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
	model.InitThirdVideo(db)

	logger.SetLevel(config.GetLoggerLevel())

	apiurl := "http://app.pearvideo.com/clt/jsp/v4/getChannels.jsp"
	GetPearCategory(apiurl,db)

	//apiurl := "http://app.pearvideo.com/clt/jsp/v4/home.jsp"
	//
	//GetPearVideoPageContent(apiurl, db)
}

type PearTag struct {
	TagId string `json:"tagId"`
	Name  string `json:"name"`
}

type PearVideo struct {
	VideoId  string `json:"videoId"`
	Url      string `json:"url"`
	FileSize string `json:"fileSize"`
	Duration string `json:"duration"`
}

type PearGeo struct {
	NamePath  string  `json:"namePath"`
	ShowName  string  `json:"showName"`
	Address   string  `json:"address"`
	Loc       string  `json:"loc"`
	PlaceName string  `json:"placeName"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type PearUser struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname"`
	Pic      string `json:"pic"`
	Level    string `json:"level"`
	IsFollow string `json:"isFollow"`
}
type PearContentList struct {
	ContId       string      `json:"contId"`
	Name         string      `json:"name"`
	Subtitle     string      `json:"subtitle"`
	Pic          string      `json:"pic"`
	UserInfo     PearUser    `json:"userInfo"`
	Link         string      `json:"link"`
	LinkType     string      `json:"linkType"`
	ForwordType  string      `json:"forwordType"`
	Duration     string      `json:"duration"`
	CommentTimes string      `json:"commentTimes"`
	Summary      string      `json:"summary"`
	Tag          []PearTag   `json:"tags"`
	Video        []PearVideo `json:"videos"`
	Geo          PearGeo     `json:"geo"`
}
type PearDataList struct {
	NodeType string            `json:"nodeType"`
	NodeName string            `json:"nodeName"`
	MoreId   string            `json:"moreId"`
	ContList []PearContentList `json:"contList"`
}
type Pear struct {
	ResultCode string         `json:"resultCode"`
	ResultMsg  string         `json:"resultMsg"`
	ReqId      string         `json:"reqId"`
	SystemTime string         `json:"systemTime"`
	NextUrl    string         `json:"nextUrl"`
	DataList   []PearDataList `json:"dataList"`
}

type PearChannelList struct {
	Type  string         `json:"type"`
	CategoryId  string         `json:"categoryId"`
	Name  string         `json:"name"`
	IsSort  string         `json:"isSort"`
}
type PearCategoryInfo struct {
	ResultCode string         `json:"resultCode"`
	ResultMsg  string         `json:"resultMsg"`
	ReqId      string         `json:"reqId"`
	SystemTime string         `json:"systemTime"`
	Channels   []PearChannelList `json:"channelList"`
}

type PearLocalChannel struct {
	ChannelCode string         `json:"channelCode"`
	Name string         `json:"name"`
	NameSpell string         `json:"nameSpell"`
	AliasName string         `json:"aliasName"`
	IsLocal string         `json:"isLocal"`
	IsLifeCircle string         `json:"isLifeCircle"`

}
type PearLocalChannelInfo struct {
	ResultCode string         `json:"resultCode"`
	ResultMsg  string         `json:"resultMsg"`
	ReqId      string         `json:"reqId"`
	SystemTime string         `json:"systemTime"`
	AutoLocalChannelInfo   []PearLocalChannel `json:"channelList"`
	HotChannelList   []PearLocalChannel `json:"channelList"`
	ChannelList   []PearLocalChannel `json:"channelList"`

}
func GetPearCategory(apiurl string,db *gorm.DB)bool{

	requ, err := http.NewRequest("GET", apiurl, nil)
	requ.Header.Add("Host", "app.pearvideo.com")
	requ.Header.Add("User-Agent", "LiVideoIOS/4.3.6 (iPhone; iOS 11.3.1; Scale/3.00)")
	requ.Header.Add("Cookie", "JSESSIONID=5F1DAADBB1FE3C7977FC373D2DEB99DA; __ads_session=IIpZt982HAmm5uKZBgA=; PEAR_PLATFORM=1; PEAR_UUID=F7835587-9CBD-4B92-8C6E-2813811E7B5F; PV_APP=srv-pv-prod-portal1")

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

	var categorys PearCategoryInfo
	if err = json.Unmarshal(recv, &categorys); err != nil {
		logger.Error(err)
		return false
	}

	if categorys.ResultMsg != "success" {
		logger.Error(errors.New("梨视频接口返回数据异常!!!!"))
		return false
	}

	for _,category := range categorys.Channels{
		logger.Debug(category.CategoryId , "   " , category.Name)
		var method string
		if category.Name == "精选"{
			apiurl = "http://app.pearvideo.com/clt/jsp/v4/home.jsp"
			method = "POST"
		}else if category.Name == "万象"{
			method = "GET"
			apiurl = "http://app.pearvideo.com/clt/jsp/v4/getNewsList.jsp"
		}else if category.Name == "当地" {
			method = "GET"
			//获取地方列表
			apiurl = "http://app.pearvideo.com/clt/jsp/v4/localChannels.jsp"
			requ, err := http.NewRequest(method, apiurl, nil)
			requ.Header.Add("Host", "app.pearvideo.com")
			requ.Header.Add("User-Agent", "LiVideoIOS/4.3.6 (iPhone; iOS 11.3.1; Scale/3.00)")
			//requ.Header.Add("Cookie", "JSESSIONID=5F1DAADBB1FE3C7977FC373D2DEB99DA; __ads_session=eG6INSc2HAm6SuGYLAA=; PEAR_PLATFORM=1; PEAR_UUID=F7835587-9CBD-4B92-8C6E-2813811E7B5F; PV_APP=srv-pv-prod-portal1")
			requ.Header.Add("Cookie", "JSESSIONID=5F1DAADBB1FE3C7977FC373D2DEB99DA; __ads_session=IIpZt982HAmm5uKZBgA=; PEAR_PLATFORM=1; PEAR_UUID=F7835587-9CBD-4B92-8C6E-2813811E7B5F; PV_APP=srv-pv-prod-portal1")

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

			var localChannels PearLocalChannelInfo
			if err = json.Unmarshal(recv, &localChannels); err != nil {
				logger.Error(err)
				return false
			}

			if localChannels.ResultMsg != "success" {
				logger.Error(errors.New("梨视频接口返回数据异常!!!!"))
				return false
			}

			for _, channel := range localChannels.ChannelList{
				apiurl = "http://app.pearvideo.com/clt/jsp/v4/localChannelConts.jsp?channelCode=" + channel.ChannelCode + "&hotPageidx=1"
				GetPearVideoPageContent(channel.Name,method,apiurl,db)
			}
			//根据地方列表获取视频信息
		}else if category.Name == "兴趣"{
			method = "POST"
			apiurl = "http://app.pearvideo.com/clt/jsp/v4/getVodConts.jsp"
		}else{
			method = "GET"
			apiurl = "http://app.pearvideo.com/clt/jsp/v4/getCategoryConts.jsp?categoryId=" + category.CategoryId + "&hotPageidx=1"
		}
		if category.Name != "当地"{
			GetPearVideoPageContent(category.Name,method,apiurl,db)
		}
	}

	return true
}

func GetPearVideoPageContent(category,method ,apiurl string, db *gorm.DB) bool {

	requ, err := http.NewRequest(method, apiurl, nil)
	requ.Header.Add("Host", "app.pearvideo.com")
	requ.Header.Add("User-Agent", "LiVideoIOS/4.3.6 (iPhone; iOS 11.3.1; Scale/3.00)")
	//requ.Header.Add("Cookie", "JSESSIONID=5F1DAADBB1FE3C7977FC373D2DEB99DA; __ads_session=eG6INSc2HAm6SuGYLAA=; PEAR_PLATFORM=1; PEAR_UUID=F7835587-9CBD-4B92-8C6E-2813811E7B5F; PV_APP=srv-pv-prod-portal1")
	requ.Header.Add("Cookie", "JSESSIONID=5F1DAADBB1FE3C7977FC373D2DEB99DA; __ads_session=IIpZt982HAmm5uKZBgA=; PEAR_PLATFORM=1; PEAR_UUID=F7835587-9CBD-4B92-8C6E-2813811E7B5F; PV_APP=srv-pv-prod-portal1")

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

	var pear Pear
	if err = json.Unmarshal(recv, &pear); err != nil {
		logger.Error(err)
		return false
	}

	if pear.ResultMsg != "success" {
		logger.Error(errors.New("梨视频接口返回数据异常!!!!"))
		return false
	}

	for i, data := range pear.DataList {
		for _, content := range data.ContList {
			var thirdVideo model.ThirdVideo
			thirdVideo.Provider = "pear"
			thirdVideo.Title = content.Name
			thirdVideo.Description = content.Summary
			thirdVideo.ThumbX = content.Pic
			for _, tag := range content.Tag {
				if tag.Name == "内容质量差" {
					continue
				} else {
					thirdVideo.Tag += tag.Name + ","
				}
			}
			if len(thirdVideo.Tag) > 0 {
				thirdVideo.Tag = thirdVideo.Tag[0 : len(thirdVideo.Tag)-1]
			}
			p := time.Now()

			for _, v := range content.Video {
				if strings.Contains(v.Url, "-hd.") {
					thirdVideo.Playurl = v.Url
					//thirdVideo.FileName = "/episode/2018/06/08/" + name
					fileSize, _ := strconv.Atoi(v.FileSize)
					thirdVideo.Filesize = uint64(fileSize)
					thirdVideo.Duration = v.Duration
					thirdVideo.ThirdVideoId = v.VideoId
					break
				}
			}

			if len(thirdVideo.Playurl) == 0 {
				continue
			}

			thirdVideo.Address = content.Geo.PlaceName
			thirdVideo.SimpleAddr = content.Geo.Address
			addressInfo := strings.Split(content.Geo.NamePath, ",")
			if len(addressInfo) == 4 {
				thirdVideo.Country = addressInfo[0]
				thirdVideo.Province = addressInfo[1]
				thirdVideo.City = addressInfo[2]
				thirdVideo.District = addressInfo[3]
			} else if len(addressInfo) == 3 {
				thirdVideo.Country = addressInfo[0]
				thirdVideo.Province = addressInfo[1]
				thirdVideo.City = addressInfo[2]
			} else if len(addressInfo) == 2 {
				thirdVideo.Country = addressInfo[0]
				thirdVideo.Province = addressInfo[1]
			} else if len(addressInfo) == 1 {
				thirdVideo.Country = addressInfo[0]
			}
			thirdVideo.Longitude = content.Geo.Longitude
			thirdVideo.Latitude = content.Geo.Latitude
			times, _ := strconv.Atoi(content.CommentTimes)
			thirdVideo.CommentCount = uint32(times)
			thirdVideo.NickName = content.UserInfo.Nickname
			thirdVideo.AuthorThumb = content.UserInfo.Pic
			thirdVideo.ThirdAuthorId = content.UserInfo.UserId
			thirdVideo.IsVerticalScreen = false
			thirdVideo.CreatedAt = p
			thirdVideo.UpdatedAt = p
			thirdVideo.Category = category

			if err := db.Where("provider = ? and third_video_id = ?", thirdVideo.Provider, thirdVideo.ThirdVideoId).First(&thirdVideo).Error; err == nil {
				continue
			}

			//PearDownloadFile(thirdVideo.Playurl, thirdVideo.FileName)

			name := fmt.Sprintf("%04d%02d%02d%02d%02d%02d_%d_1200_1.MP4", p.Year(), p.Month(), p.Day(), p.Hour(), p.Minute(), p.Second(), (i + 1))
			thirdVideo.FileName = "/episode/" + fmt.Sprint(p.Year()) + "/" + fmt.Sprintf("%02d",p.Month()) + "/" + fmt.Sprintf("%02d",p.Day()) + "/" + name

			if err = db.Create(&thirdVideo).Error; err != nil {
				logger.Error(err)
				return false
			}
		}

	}

	if len(pear.NextUrl) > 0 {
		GetPearVideoPageContent(category,method,pear.NextUrl, db)
	}
	return true
}

func PearDownloadFile(requrl string, filename string) (int64, error) {
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

	return bytes, nil
}
