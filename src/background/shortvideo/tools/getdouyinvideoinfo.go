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

		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822607&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822607237&as=a195bb518ff09b79605185&cp=b400b858f709159fe1ptfm&mas=00fe7ff836caff7598244eb245225a12d4ac1c8cac0c6c9cec469c",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822614&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822613669&as=a1557b81f651db29f04490&cp=bf1eb35a6b0d1893e1klbv&mas=004c80ab683856b3cdebe98a53c53fed0a0c9c2c2c0c669c4c4626",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822620&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822619132&as=a1b5fb81cc41cb69009571&cp=b61ab85dc20e1796e1ujxi&mas=006f0a0bad0254ef9cc5cdbe45ce0bd3b88cecac9c0c0c9c6c46c6",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822627&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822627241&as=a1e56b4103b2bb29802868&cp=b92ab05b34081e9de1obhk&mas=009f9dba3e12bf988ce9299219989fb4f81c6c1c4c0c1c9c4c4646",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822634&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822634144&as=a1a53b511a52fb79503283&cp=bf28bd58ad041e98e1qjpv&mas=003e73e93f5daacfc95c67dd69f16d9196cc1c4ccc0cac9cec4666",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822636&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822635789&as=a1b57bd1ccd2fb29203686&cp=b92cba54cb0f1c9ee1ayfv&mas=008f5a0bada381c88d248730b98f59b2286c1c6ccc0c4c9c4c4666",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822641&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822641926&as=a1855b5191e34b49e07765&cp=b434b6591906109be1yawd&mas=00d5678f77410f1b0febd44ed31cd2a116ac6cecec0ca69c2c462c",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822648&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822647746&as=a1655b8148732b49700452&cp=b33cba508d041c92e1fmkn&mas=002037c5fa58c1bc0f6b4e130ba3460e924cac2c0c0cec9c2c464c",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822649&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822648679&as=a1d5ab112923db79f01454&cp=b53bb55499061b95e1akdp&mas=006bcb107cfbde0f9baedb44c586355b002cac2c8c0c669cec4626",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822661&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822661080&as=a1157b6125a4bb99704514&cp=b04fbb5b5b0e1e9be1uado&mas=00e833a93ce5d6d88da0c96e67cb5ca2e22c8cac2c0cec9c9c4646",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822663&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822662566&as=a1a5fb7127a41bc9e05190&cp=b645b8517b021b9de1cghz&mas=00ab00bf3c5d3918dc8bba59692f104b560c9c8cac0ca69cc6468c",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822672&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822672383&as=a1c58b01e0b56b89104480&cp=b350b85a05091f98e1lrhj&mas=00125b78ebc5a7a3c2c275c979ea6851320c1c2c2c0c8c9c1c466c",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822674&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822673863&as=a1d50bc1d205ab09b05288&cp=b950bc5e28081c94e1chmh&mas=00fda21b2a711f8aca1c2448c38c0c8c061c1c4cac0c469c0c4686",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822678&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822679657&as=a125abc146c53b09508352&cp=b152bc5067081990e1yfqs&mas=00fb31cafa386c359b4d7f5701bce6710a4caccc1c0cac9c0c46cc",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822682&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822683631&as=a1c58b012a854bb9507336&cp=b55cbd5ba4031793e1qrkw&mas=00468129bcfb33f3c29a85fb1191044f206cccccec0cac9c46462c",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822692&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822691973&as=a1358b21e486cb89402475&cp=b66ab758460c1e97e1ezpk&mas=00a9b8e8eba2751e82e2002c6b4cf7f8ecacec2c4c0c2c9c1c46c6",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822694&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822694219&as=a1e51b11c6669be9d06335&cp=b863b8536a041f93e1dsri&mas=00164c3f6d2a0b540b27a2ff790e09b748accccc6c0c269ca6469c",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822702&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822701821&as=a1f5cb71ee26db99b00921&cp=b06ab859e0041793e1xjri&mas=00811a0cabf26f90d93d42e7fb69cf3dc48c4c9c0c0c469c9c4626",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822705&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822705077&as=a1153bb191b7ab59700640&cp=b878b6581f07129de1rdwr&mas=0090ecad3752490889310dc85fb12520b00c2c6c0c0cec9cac4686",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822714&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822714022&as=a1c5bb710a47ebc9c08806&cp=b87ab75ead081491e1ckcw&mas=0044d3a9be09639398512561535536d9de6c0c1c1c0cc69cc646a6",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822721&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822721029&as=a1d5cbe101088bd9c07995&cp=b387b55f150a1a91e1jhzy&mas=00cf6d187e991cc4d9c5af85f16d6637faac9c9cec0cc69c26461c",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822734&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822733553&as=a1c5fbf15e585b29401765&cp=b285b950e00f179be1rpbm&mas=0024d3adbb717bbc5ceafdfa4b169176a6ac6cec8c0c2c9c4c46ac",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822736&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822735268&as=a1e52ba1d0e91b79607167&cp=b09eb35c0d01179ee1vonq&mas=009ac8c96a833dee889ce2bd5d88338f38ec6c8cec0c6c9cec468c",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822743&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822743208&as=a1b57bf1f7d9db19307728&cp=bb92b65574061c99e1uouw&mas=000e69dfa8dc89d40d509fc199d911b3561c4cecec0ccc9c8c4626",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822753&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822753064&as=a1f59bf121eadb69706240&cp=bda0b35f19081e92e1morv&mas=0063a458fcf65f1bc30ee26945dfae15040c2c4c6c0cec9c6c4626",
		"https://api.amemv.com/aweme/v1/feed/?type=0&max_cursor=0&min_cursor=-1&count=6&volume=0.0&pull_type=2&need_relieve_aweme=0&ts=1527822758&app_type=normal&os_api=23&device_type=Letv%20X501&device_platform=android&ssmix=a&iid=30952143538&manifest_version_code=179&dpi=420&uuid=868897026417337&version_code=179&app_name=aweme&version_name=1.7.9&openudid=5a607035cb966df2&device_id=51033988626&resolution=1080*1920&os_version=6.0&language=zh&device_brand=Letv&ac=wifi&update_version_code=1792&aid=1128&channel=yyh_search&_rticket=1527822757474&as=a1051b2166ea4bf9d08347&cp=b8aeba5064011392e1qscn&mas=001d088a78584f840baec1717f9cc47eb2ec2ccc1c0c269c66462c",
	}

	for k, apiurl := range apiurls {
		requ, err := http.NewRequest("GET", apiurl, nil)
		requ.Header.Add("Host", "api.amemv.com")
		requ.Header.Add("User-Agent", "okhttp/3.8.1")

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

			duration := uint32(aweme.Video.Duration / 1000)

			video.Description = strings.Replace(video.Description, "抖音", "", -1)
			now := time.Now()
			video.CreatedAt = now
			video.UpdatedAt = now
			video.Duration = duration
			video.ThumbY = aweme.Video.Cover.UrlList[0]
			video.Duration = duration
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
				return
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
