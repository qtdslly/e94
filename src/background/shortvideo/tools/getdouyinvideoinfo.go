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
	logger.SetLevel(config.GetLoggerLevel())

	for {
		if !GetDouYinPageContent(db) {
			break
		}
	}
}

type DouYinRiskInfos struct {
	Warn     bool   `json:"warn"`
	Content  string `json:"content"`
	RiskSink bool   `json:"risk_sink"`
	Type     uint8  `json:"type"`
}

type DouYinIcon struct {
	UrlList []string `gorm:"url_list" json:"url_list"`
	Uri     string   `gorm:"uri" json:"uri"`
}

type DouYinStatistics struct {
	PlayCount    uint32 `json:"play_count"`
	AwemeId      string `json:"aweme_id"`
	CommentCount uint32 `json:"comment_count"`
	ShareCount   uint32 `json:"share_count"`
	DiggCount    uint32 `json:"digg_count"`
}

type DouYinStatus struct {
	WithGoods       bool  `json:"with_goods"`
	IsDelete        bool  `json:"is_delete"`
	PrivateStatus   uint8 `json:"private_status"`
	WithFusionGoods bool  `json:"with_fusion_goods"`
	AllowComment    bool  `json:"allow_comment"`
	AllowShare      bool  `json:"allow_share"`
	IsPrivate       bool  `json:"is_private"`
}

type DouYinShareInfo struct {
	ShareWeiboDesc     string `json:"share_weibo_desc"`
	BoolPersist        uint8  `json:"bool_persist"`
	ShareQuote         string `json:"share_quote"`
	ShareTitle         string `json:"share_title"`
	ShareSignatureDesc string `json:"share_signature_desc"`
	ShareSignatureUrl  string `json:"share_signature_url"`
	ShareUrl           string `json:"share_url"`
	ShareDesc          string `json:"share_desc"`
}

type DouYinVideo struct {
	Ratio       string     `json:"ratio"`
	OriginCover DouYinIcon `json:"origin_cover"`
	PlayAddr    DouYinIcon `json:"play_addr"`
	Cover       DouYinIcon `json:"cover"`
	Height      uint64     `json:"height"`
	Width       uint64     `json:"width"`
	BitRate     []*struct {
		BitRate     uint32 `json:"bit_rate"`
		GearName    string `json:"gear_name"`
		QualityType uint8  `json:"quality_type"`
	} `json:"bit_rate"`
	DownloadAddr  DouYinIcon `json:"download_addr"`
	HasWatermark  bool       `json:"has_watermark"`
	Duration      uint32     `json:"duration"`
	PlayAddrLowbr DouYinIcon `json:"play_addr_lowbr"`
	DynamicCover  DouYinIcon `json:"dynamic_cover"`
}

type DouYinAutor struct {
	AuthorityStatus        uint8      `json:"authority_status"`
	WeiboSchema            string     `json:"weibo_schema"`
	Uid                    string     `json:"uid"`
	BindPhone              string     `json:"bind_phone"`
	WeiboUrl               string     `json:"weibo_url"`
	ShareQrcodeUri         string     `json:"share_qrcode_uri"`
	LiveAgreementTime      uint64     `json:"live_agreement_time"`
	YoutubeChannelId       string     `json:"youtube_channel_id"`
	WeiboName              string     `json:"weibo_name"`
	VideoIcon              DouYinIcon `json:"video_icon"`
	CommerceUserLevel      uint8      `json:"commerce_user_level"`
	AvatarLarger           DouYinIcon `json:"avatar_larger"`
	CreateTime             uint64     `json:"create_time"`
	LiveAgreement          uint8      `json:"live_agreement"`
	EnterpriseVerifyReason string     `json:"enterprise_verify_reason"`
	OriginalMusicQrcode    string     `json:"original_music_qrcode"`
	IsAdFake               bool       `json:"is_ad_fake"`
	VerifyInfo             string     `json:"verify_info"`
	GoogleAccount          string     `json:"google_account"`
	Constellation          uint64     `json:"constellation"`
	YoutubeChannelTitle    string     `json:"youtube_channel_title"`
	StoryOpen              string     `json:"story_open"`
	SchoolType             uint8      `json:"school_type"`
	DuetSetting            uint8      `json:"duet_setting"`
	AcceptPrivatePolicy    bool       `json:"accept_private_policy"`
	NeedRecommend          uint8      `json:"need_recommend"`
	ShortId                string     `json:"short_id"`
	TwitterId              string     `json:"twitter_id"`
	PreventDownload        bool       `json:"prevent_download"`
	VerificationType       uint8      `json:"verification_type"`
	FollowerStatus         uint8      `json:"follower_status"`
	AccountRegion          string     `json:"account_region"`
	RoomId                 uint64     `json:"room_id"`
	AvatarMedium           DouYinIcon `json:"avatar_medium"`
	ReflowPageGid          uint64     `json:"reflow_page_gid"`
	AvatarThumb            DouYinIcon `json:"story_open"`
	IsBindedWeibo          bool       `json:"is_binded_weibo"`
	IsVerified             bool       `json:"is_verified"`
	VideoIconVirtualUri    string     `json:"video_icon_virtual_URI"`
	TwitterName            string     `json:"twitter_name"`
	HideSearch             bool       `json:"hide_search"`
	UserCanceled           bool       `json:"user_canceled"`
	OriginalMusicCover     DouYinIcon `json:"original_music_cover"`
	WithCommerceEntry      bool       `json:"with_commerce_entry"`
	AppleAccount           uint64     `json:"apple_account"`
	FollowStatus           uint8      `json:"follow_status"`
	SchoolName             string     `json:"school_name"`
	HasOrders              bool       `json:"has_orders"`
	Birthday               string     `json:"birthday"`
	CustomVerify           string     `json:"custom_verify"`
	InsId                  string     `json:"ins_id"`
	UserRate               string     `json:"user_rate"`
	NickName               string     `json:"nickname"`
	SpecialLock            uint8      `json:"special_lock"`
	CommentSetting         uint8      `json:"comment_setting"`
	LiveVerify             uint8      `json:"live_verify"`
	ShieldFollowNotice     uint8      `json:"shield_follow_notice"`
	ShieldCommentNotice    uint8      `json:"shield_comment_notice"`
	ShieldDiggNotice       uint8      `json:"shield_digg_notice"`
	HideLocation           bool       `json:"hide_location"`
	Gender                 uint8      `json:"gender"`
	Region                 string     `json:"region"`
	Secret                 string     `json:"custom_verify"`
	UniqueIdModifyTime     uint64     `json:"unique_id_modify_time"`
	SchoolPoiDId           string     `json:"school_poi_id"`
	Status                 uint8      `json:"status"`
	PolicyVersion          *struct {
	} `json:"policy_version"`
	AvatarUri     string   `json:"avatar_uri"`
	Signature     string   `json:"signature"`
	Geofencing    []string `json:"geofencing"`
	WeiboVerify   string   `json:"weibo_verify"`
	IsPhoneBinded bool     `json:"is_phone_binded"`
	UniqueId      string   `json:"unique_id"`
}

type DouYinAddressInfo struct {
	Province   string `json:"province"`
	City       string `json:"city"`
	SimpleAddr string `json:"simple_addr"`
	District   string `json:"district"`
	Address    string `json:"address"`
}
type DouYinPoiInfo struct {
	AddressInfo DouYinAddressInfo `json:"address_info"`
}

type DouYinAwemeList struct {
	RiskInfos           DouYinRiskInfos `json:"risk_infos"`
	LabelTop            DouYinIcon      `json:"label_top"`
	AuthorUserId        uint64          `json:"author_user_id"`
	ItemCommentSettings uint8           `json:"item_comment_settings"`
	Rate                int32           `json:"rate"`
	CreateTime          int64           `json:"create_time"`
	Video               DouYinVideo     `json:"video"`
	AwemeId             string          `json:"aweme_id"`
	VideoLabels         []*struct {
	} `json:"video_labels"`
	IsVr           bool             `json:"is_vr"`
	VrType         uint8            `json:"vr_type"`
	Statistics     DouYinStatistics `json:"statistics"`
	Author         DouYinAutor      `json:"author"`
	CmtSwt         bool             `json:"cmt_swt"`
	ShareUrl       string           `json:"share_url"`
	IsAd           bool             `json:"is_ad"`
	Music          string           `json:"music1"`
	BodydanceScore uint32           `json:"bodydance_score"`
	IsHashTag      uint8            `json:"is_hash_tag"`
	Status         DouYinStatus     `json:"status"`
	SortLabel      string           `json:"sort_label"`
	Descendants    *struct {
		Platforms []string `json:"platforms"`
		NotifyMsg string   `json:"notify_msg"`
	} `json:"descendants"`
	ShareInfo   DouYinShareInfo `json:"share_info"`
	VideoText   []string        `json:"video_text"`
	IsTop       uint8           `json:"is_top"`
	CollectStat uint8           `json:"collect_stat"`
	AwemeType   uint8           `json:"aweme_type"`
	Desc        string          `json:"desc"`
	PoiInfo     DouYinPoiInfo   `json:"poi_info"`
	Geofencing  []string        `json:"geofencing"`
	Region      string          `json:"region"`
	IsPgcshow   bool            `json:"is_pgcshow"`
	IsRelieve   bool            `json:"is_relieve"`
	TextExtra   []*struct {
		Start  uint32 `json:"start"`
		UserId string `json:"user_id"`
		End    uint32 `json:"end"`
		Type   uint8  `json:"type"`
	} `json:"text_extra"`
	UserDigged uint8 `json:"user_digged"`
}

type DouYinExtra struct {
	Logid        string   `json:"logid"`
	Now          uint64   `json:"now"`
	FatalItemIds []string `json:"fatal_item_ids"`
}
type DouYinData struct {
	MaxCursor    uint              `json:"max_cursor"`
	AwemeList    []DouYinAwemeList `json:"aweme_list"`
	RefreshClear uint              `json:"refresh_clear"`
	Extra        DouYinExtra       `json:"extra"`
	HasMore      uint              `json:"has_more"`
	StatusCode   uint              `json:"status_code"`
	HomeModel    uint              `json:"home_model"`
	Rid          string            `json:"rid"`
	MinCursor    int               `json:"min_cursor"`
}

func GetDouYinPageContent(db *gorm.DB) bool {
	tx := db.Begin()

	//params := make(map[string]interface{})
	//params["type"] = "0"
	//params["max_cursor"] = "0"
	//params["min_cursor"] = "-1"
	//params["count"] = "6"
	//params["volume"] = "0.0"
	//params["pull_type"] = "2"   //0
	//params["need_relieve_aweme"] = "0"
	//params["app_type"] = "normal"
	//params["os_api"] = "23"
	//params["device_platform"] = "android"
	//params["ssmix"] = "a"
	//params["iid"] = "30952143538"
	//params["manifest_version_code"] = "179"
	//params["dpi"] = "420"
	//params["uuid"] = "868897026417337"
	//params["version_code"] = "179"
	//params["app_name"] = "aweme"
	//params["version_name"] = "1.7.9"
	//params["openudid"] = "5a607035cb966df2"
	//params["device_id"] = "51033988626"
	//params["os_version"] = "6.0"
	//params["language"] = "zh"
	//params["device_brand"] = "Letv"
	//params["ac"] = "wifi"
	//params["resolution"] = "1080*1920"
	//params["device_type"] = "Letv X501"
	//params["update_version_code"] = "1792"
	//params["aid"] = "1128"
	//params["channel"] = "yyh_search"
	//params["ts"] = "1527756780"
	//params["_rticket"] = "1527756779246"
	//params["as"] = "a1758bd03c3e3bd79f9905"
	//params["cp"] = "bae6b15dccf20a74e1tgnt"
	//params["mas"] = "00be88aebd1e19608264246d4d980b9898ac0c9c9c669cec2646cc"

	//values := url.Values{}
	//for k, v := range params {
	//	values.Add(k, fmt.Sprint(v))
	//}
	//apiurl := "https://api.amemv.com/aweme/v1/feed/?" + values.Encode()
	//var apiurl string
	//apiurl = "https://aweme.snssdk.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527757160&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527757158963&as=a1a55bd018463b295f9547&cp=b664b95482f40392e1wlfe&mas=0040ff203fc84ed3cfba094c53503af156ec2cac9c66ac9c4c46cc"

	apiurls := []string{
		"https://aweme-eagle-hl.snssdk.com/aweme/v1/feed/?version_code=5.8.0&pass-region=1&pass-route=1&js_sdk_version=1.13.3.0&app_name=aweme&vid=505ECF2B-061B-40AF-AA0E-2F4DB573FCAA&app_version=5.8.0&device_id=51632502948&channel=App%20Store&mcc_mnc=46001&aid=1128&screen_width=1242&openudid=0912a38b4938c6d19a0e2f0a745339deb70bbe82&os_api=18&ac=WIFI&os_version=12.2&device_platform=iphone&build_number=58010&device_type=iPhone8,2&iid=70375126620&idfa=F7835587-9CBD-4B92-8C6E-2813811E7B5F&volume=0.00&count=6&longitude=114.3404602033706&feed_style=0&filter_warn=0&address_book_access=1&user_id=102900047091&type=0&latitude=30.55565496650956&gps_access=3&pull_type=2&max_cursor=0&mas=01775e4a637fafd29a0e52c72f5763221a8856dae5bf755107eeb7&as=a2d5cc4cf7161cc9222904&ts=1556269415",
}

	for k, apiurl := range apiurls {
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

		logger.Debug(string(recv))

		var douyin DouYinData
		if err = json.Unmarshal(recv, &douyin); err != nil {
			logger.Error(err)
			return false
		}

		if douyin.StatusCode != 0 {
			logger.Error(errors.New("抖音接口返回数据异常!!!!"))
			if k == len(apiurls)-1 {
				return false
			}
			continue
		}

		for i, aweme := range douyin.AwemeList {
			if aweme.IsAd {
				continue
			}
			logger.Debug("video_id:", aweme.AwemeId)

			var video model.ThirdVideo
			video.Provider = "douyin"
			video.ThirdVideoId = aweme.AwemeId

			if err = db.Where("provider = ? and third_video_id = ?", video.Provider, video.ThirdVideoId).First(&video).Error; err == nil {
				continue
			}

			p := time.Now()
			name := fmt.Sprintf("%04d%02d%02d%02d%02d%02d_%d_1000_1.MP4", p.Year(), p.Month(), p.Day(), p.Hour(), p.Minute(), p.Second(), (i + 1))

			if len(aweme.Video.PlayAddr.UrlList) > 0 {
				DouYinDownloadFile(aweme.Video.PlayAddr.UrlList[0], name)
			} else {
				continue
			}

			logger.Debug("url:", video.Playurl)
			logger.Debug("author_id:", aweme.Author.Uid)
			logger.Debug("author_short_id:", aweme.Author.ShortId)
			logger.Debug("author_brithday:", aweme.Author.Birthday)
			logger.Debug("nick_name:", aweme.Author.NickName)
			logger.Debug("video_description:", aweme.Desc)
			logger.Debug("play_url:", aweme.Video.PlayAddr.UrlList[0])
			logger.Debug("thumb:", aweme.Video.DynamicCover.UrlList[0])
			logger.Debug("width:", aweme.Video.Width)
			logger.Debug("height:", aweme.Video.Height)
			logger.Debug("duration:", aweme.Video.Duration)
			logger.Debug("has_water_mark:", aweme.Video.HasWatermark)
			logger.Debug("play_count:", aweme.Statistics.PlayCount)
			logger.Debug("share_count:", aweme.Statistics.ShareCount)
			logger.Debug("comment_count:", aweme.Statistics.CommentCount)
			logger.Debug("digg_count:", aweme.Statistics.DiggCount)
			logger.Debug("author_thumb:", aweme.Author.AvatarLarger.UrlList[0])
			logger.Debug("share_title:", aweme.ShareInfo.ShareTitle)
			logger.Debug("share_description:", aweme.ShareInfo.ShareDesc)

			//duration := uint32(aweme.Video.Duration / 1000)

			video.Description = strings.Replace(video.Description, "抖音", "", -1)
			now := time.Now()
			video.CreatedAt = now
			video.UpdatedAt = now
			//video.Duration = duration
			video.ThumbY = aweme.Video.Cover.UrlList[0]
			//video.Duration = duration
			video.Width = uint32(aweme.Video.Width)
			video.Height = uint32(aweme.Video.Height)
			video.Playurl = "/episode/2018/06/11/" + name
			video.DiggCount = aweme.Statistics.DiggCount
			video.PlayCount = aweme.Statistics.PlayCount
			video.CommentCount = aweme.Statistics.CommentCount
			video.ShareCount = aweme.Statistics.ShareCount
			video.Province = aweme.PoiInfo.AddressInfo.Province
			video.City = aweme.PoiInfo.AddressInfo.City
			video.District = aweme.PoiInfo.AddressInfo.District
			video.Address = aweme.PoiInfo.AddressInfo.Address

			video.ThirdAuthorId = aweme.Author.Uid
			video.NickName = aweme.Author.NickName
			if len(aweme.Author.AvatarThumb.UrlList) > 0{
				video.AuthorThumb = aweme.Author.AvatarThumb.UrlList[0]
			}
			video.Birthday = aweme.Author.Birthday
			video.IsVerticalScreen = true
			video.HasWaterMark = aweme.Video.HasWatermark

			if err = tx.Create(&video).Error; err != nil {
				logger.Error(err)
				return false
			}
		}
	}

	return true
}

func DouYinDownloadFile(requrl string, filename string) (int64, error) {
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
