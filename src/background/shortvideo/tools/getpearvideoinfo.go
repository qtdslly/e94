package main

import (
	"cms/config"
	"cms/model"
	"cms/setting"
	"common/constant"
	"common/logger"

	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"net/http"
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
	model.InitModel(db)

	logger.SetLevel(config.GetLoggerLevel())

	var property model.Property
	if err := db.Where("name = '类型'").First(&property).Error; err != nil {
		logger.Error(err)
		return
	}

	resAdrr, err := setting.GetResBindAddr(db)
	if err != nil {
		logger.Error(err)
		return
	}

	apiurl := "http://app.pearvideo.com/clt/jsp/v4/getChannels.jsp"
	GetPearCategory(resAdrr, property.Id, apiurl, db)
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
type PearContent struct {
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
	NodeType string        `json:"nodeType"`
	NodeName string        `json:"nodeName"`
	MoreId   string        `json:"moreId"`
	ContList []PearContent `json:"contList"`
}
type PearJinXuan struct {
	ResultCode string         `json:"resultCode"`
	ResultMsg  string         `json:"resultMsg"`
	ReqId      string         `json:"reqId"`
	SystemTime string         `json:"systemTime"`
	NextUrl    string         `json:"nextUrl"`
	DataList   []PearDataList `json:"dataList"`
}

type PearChannelList struct {
	Type       string `json:"type"`
	CategoryId string `json:"categoryId"`
	Name       string `json:"name"`
	IsSort     string `json:"isSort"`
}
type PearCategoryInfo struct {
	ResultCode string            `json:"resultCode"`
	ResultMsg  string            `json:"resultMsg"`
	ReqId      string            `json:"reqId"`
	SystemTime string            `json:"systemTime"`
	Channels   []PearChannelList `json:"channelList"`
}

type PearLocalChannel struct {
	ChannelCode  string `json:"channelCode"`
	Name         string `json:"name"`
	NameSpell    string `json:"nameSpell"`
	AliasName    string `json:"aliasName"`
	IsLocal      string `json:"isLocal"`
	IsLifeCircle string `json:"isLifeCircle"`
}
type PearLocalChannelInfo struct {
	ResultCode string `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
	ReqId      string `json:"reqId"`
	SystemTime string `json:"systemTime"`
	//AutoLocalChannelInfo   []PearLocalChannel `json:"autoLocalChannelInfo"`
	HotChannelList []PearLocalChannel `json:"hotChannelList"`
	ChannelList    []PearLocalChannel `json:"channelList"`
}

type PearLocal struct {
	ResultCode string        `json:"resultCode"`
	ResultMsg  string        `json:"resultMsg"`
	ReqId      string        `json:"reqId"`
	SystemTime string        `json:"systemTime"`
	NextUrl    string        `json:"nextUrl"`
	ContList   []PearContent `json:"contList"`
}

type PearOnlyContent struct {
	ContId      string      `json:"contId"`
	Name        string      `json:"name"`
	Pic         string      `json:"pic"`
	UserInfo    PearUser    `json:"userInfo"`
	Link        string      `json:"link"`
	LinkType    string      `json:"linkType"`
	ForwordType string      `json:"forwordType"`
	Duration    string      `json:"duration"`
	Tag         []PearTag   `json:"tags"`
	Video       []PearVideo `json:"videos"`
}

type PearOthers struct {
	ResultCode string             `json:"resultCode"`
	ResultMsg  string             `json:"resultMsg"`
	ReqId      string             `json:"reqId"`
	SystemTime string             `json:"systemTime"`
	NextUrl    string             `json:"nextUrl"`
	ContList   []PearOtherContent `json:"contList"`
}

type PearOtherContent struct {
	ContId       string      `json:"contId"`
	Name         string      `json:"name"`
	Pic          string      `json:"pic"`
	UserInfo     PearUser    `json:"userInfo"`
	Duration     string      `json:"duration"`
	CommentTimes string      `json:"commentTimes"`
	Summary      string      `json:"summary"`
	Tag          []PearTag   `json:"tags"`
	Video        []PearVideo `json:"videos"`
	Geo          PearGeo     `json:"geo"`
}

type PearOtherContentSecond struct {
	ResultCode string      `json:"resultCode"`
	ResultMsg  string      `json:"resultMsg"`
	ReqId      string      `json:"reqId"`
	SystemTime string      `json:"systemTime"`
	Content    PearContent `json:"content"`
}

func GetPearCategory(resAdrr string, propertyId uint32, apiurl string, db *gorm.DB) bool {

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

	for _, category := range categorys.Channels {
		logger.Debug(category.CategoryId, "   ", category.Name)
		var method string
		if category.Name == "精选" {
			apiurl = "http://app.pearvideo.com/clt/jsp/v4/home.jsp"
			method = "POST"
		} else if category.Name == "万象" {
			method = "GET"
			apiurl = "http://app.pearvideo.com/clt/jsp/v4/getNewsList.jsp"
			continue
		} else if category.Name == "当地" {
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

			logger.Debug(string(recv))

			var localChannels PearLocalChannelInfo
			if err = json.Unmarshal(recv, &localChannels); err != nil {
				logger.Error(err)
				return false
			}

			if localChannels.ResultMsg != "success" {
				logger.Error(errors.New("梨视频接口返回数据异常!!!!"))
				return false
			}

			for _, channel := range localChannels.ChannelList {
				apiurl = "http://app.pearvideo.com/clt/jsp/v4/localChannelConts.jsp?channelCode=" + channel.ChannelCode + "&hotPageidx=1"
				GetPearVideoPageContent(resAdrr, propertyId, "当地|"+channel.Name, method, apiurl, db)
			}
			//根据地方列表获取视频信息
		} else if category.Name == "兴趣" {
			method = "POST"
			apiurl = "http://app.pearvideo.com/clt/jsp/v4/getVodConts.jsp"
		} else {
			method = "GET"
			apiurl = "http://app.pearvideo.com/clt/jsp/v4/getCategoryConts.jsp?categoryId=" + category.CategoryId + "&hotPageidx=1"
		}
		if category.Name != "当地" {
			GetPearVideoPageContent(resAdrr, propertyId, category.Name, method, apiurl, db)
		}
	}

	return true
}

func GetPearVideoPageContent(resAdrr string, propertyId uint32, category, method, apiurl string, db *gorm.DB) bool {

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

	logger.Debug(string(recv))

	if category == "精选" || category == "兴趣" {
		var pear PearJinXuan
		if err = json.Unmarshal(recv, &pear); err != nil {
			logger.Error(err)
			return false
		}

		if pear.ResultMsg != "success" {
			logger.Error(errors.New("梨视频接口返回数据异常!!!!"))
			return false
		}

		for _, data := range pear.DataList {
			for _, content := range data.ContList {
				SavePearInfo(resAdrr, propertyId, content, db)
			}

		}
		if len(pear.NextUrl) > 0 {
			GetPearVideoPageContent(resAdrr, propertyId, category, method, pear.NextUrl, db)
		}
	} else if strings.Contains(category, "当地|") {
		category = strings.Replace(category, "当地|", "", -1)
		var pear PearLocal
		if err = json.Unmarshal(recv, &pear); err != nil {
			logger.Error(err)
			return false
		}

		if pear.ResultMsg != "success" {
			logger.Error(errors.New("梨视频接口返回数据异常!!!!"))
			return false
		}

		for _, content := range pear.ContList {
			SavePearInfo(resAdrr, propertyId, content, db)
		}

		if len(pear.NextUrl) > 0 {
			GetPearVideoPageContent(resAdrr, propertyId, category, method, pear.NextUrl, db)
		}
	} else {
		var pear PearOthers
		if err = json.Unmarshal(recv, &pear); err != nil {
			logger.Error(err)
			return false
		}

		if pear.ResultMsg != "success" {
			logger.Error(errors.New("梨视频接口返回数据异常!!!!"))
			return false
		}

		for _, content := range pear.ContList {
			method = "POST"
			apiurl = "http://app.pearvideo.com/clt/jsp/v4/content.jsp?contId=" + content.ContId
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

			logger.Debug(string(recv))

			var pear PearOtherContentSecond
			if err = json.Unmarshal(recv, &pear); err != nil {
				logger.Error(err)
				return false
			}

			if pear.ResultMsg != "success" {
				logger.Error(errors.New("梨视频接口返回数据异常!!!!"))
				return false
			}

			SavePearInfo(resAdrr, propertyId, pear.Content, db)
		}

		if len(pear.NextUrl) > 0 {
			GetPearVideoPageContent(resAdrr, propertyId, category, method, pear.NextUrl, db)
		}
	}

	return true
}

func SavePearInfo(resAdrr string, propertyId uint32, content PearContent, db *gorm.DB) {
	var err error
	var shortVideo model.Video
	shortVideo.Provider = constant.ContentProviderPear

	for _, v := range content.Video {
		if strings.Contains(v.Url, "-hd.") {
			shortVideo.Url = v.Url
			fileSize, _ := strconv.Atoi(v.FileSize)
			shortVideo.Filesize = uint32(fileSize / 1024 / 1024)
			duration, _ := strconv.Atoi(v.Duration)
			shortVideo.Duration = uint32(duration)
			shortVideo.SourceId = v.VideoId
			break
		}
	}

	if len(shortVideo.Url) == 0 {
		return
	}

	if err := db.Where("provider = ? and source_id = ?", shortVideo.Provider, shortVideo.SourceId).First(&shortVideo).Error; err == nil {
		return
	}

	shortVideo.Title = content.Name
	shortVideo.Description = content.Summary
	shortVideo.ThumbX = content.Pic

	shortVideo.Address = content.Geo.PlaceName
	addressInfo := strings.Split(content.Geo.NamePath, ",")
	if len(addressInfo) == 4 {
		shortVideo.Country = addressInfo[0]
		shortVideo.Province = addressInfo[1]
		shortVideo.City = addressInfo[2]
		shortVideo.District = addressInfo[3]
	} else if len(addressInfo) == 3 {
		shortVideo.Country = addressInfo[0]
		shortVideo.Province = addressInfo[1]
		shortVideo.City = addressInfo[2]
	} else if len(addressInfo) == 2 {
		shortVideo.Country = addressInfo[0]
		shortVideo.Province = addressInfo[1]
	} else if len(addressInfo) == 1 {
		shortVideo.Country = addressInfo[0]
	}
	shortVideo.Longitude = content.Geo.Longitude
	shortVideo.Latitude = content.Geo.Latitude
	shortVideo.Vertical = false

	now := time.Now()
	shortVideo.CreatedAt = now
	shortVideo.UpdatedAt = now

	tx := db.Begin()

	var person model.Person
	person.Name = content.UserInfo.Nickname
	if err = db.Where("provider_id = ? and name = ?", constant.ContentProviderPear, person.Name).First(&person).Error; err == gorm.ErrRecordNotFound {
		person.Country = shortVideo.Country
		person.Description = ""
		person.Figure = content.UserInfo.Pic
		person.Nickname = content.UserInfo.Nickname
		person.Role = model.PersonRoleTypeUper
		person.CreatedAt = now
		person.UpdatedAt = now

		if err = tx.Create(&person).Error; err != nil {
			logger.Error(err)
			tx.Rollback()
			return
		}
	}

	shortVideo.PersonId = person.Id
	if err := db.Create(&shortVideo).Error; err != nil {
		logger.Error(err)
		tx.Rollback()
		return
	}

	for _, t := range content.Tag {
		if t.Name == "内容质量差" {
			continue
		}

		var tag model.Tag
		tag.Name = t.Name
		if err := tx.Where("name = ? and property_id = ?", tag.Name, propertyId).First(&tag).Error; err == gorm.ErrRecordNotFound {
			tag.PropertyId = propertyId
			tag.Sort = 0
			now := time.Now()
			tag.CreatedAt = now
			tag.UpdatedAt = now
			if err = tx.Create(&tag).Error; err != nil {
				logger.Error(err)
				tx.Rollback()
				return
			}
		}

		if err := db.Exec("insert into video_tag(video_id,tag_id) values(?,?)", shortVideo.Id, tag.Id).Error; err != nil {
			logger.Error(err)
			tx.Rollback()
			return
		}
	}

	tx.Commit()
}
