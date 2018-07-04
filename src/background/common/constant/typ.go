package constant

// global media content type definition.
const (
	MediaTypeUnknown       = 0
	MediaTypeVideo         = 1
	MediaTypeEpisode       = 2
	MediaTypePlayUrl       = 3
	MediaTypeStream        = 4
	MediaTypeProgram       = 5
	MediaTypeEpg           = 6
	MediaTypeEpgBlock      = 7
	MediaTypeNotification  = 8
	MediaTypeBroadcast     = 9
	MediaTypeAd            = 10
	MediaTypePage          = 11
	MediaTypeSection       = 12
	MediaTypeFeed          = 13
	MediaTypeH5            = 14
	MediaTypeRes           = 15 // 资源文件，图片、apk等文件
	MediaTypeAdUnit        = 16
	MediaTypePerson        = 17
	MediaTypeResourceGroup = 18
	MediaTypeArticle       = 19
	MediaTypeGallery       = 20
	MediaTypeTimeBlock     = 21
	MediaTypeAlbum         = 22
	MediaTypePlaylist      = 23
	MediaTypeMusic         = 24
	MediaTypeApp           = 25
	MediaTypeBook          = 26
)

// global os type definition.
const (
	OsTypeUnknown = 0 // 未知类型
	OsTypeAndroid = 1 // android
	OsTypeIos     = 2 // ios
)

// 自动检测视频质量，选择最接近的分辨率
const (
	VideoQuality360p  = 1 // 480x360，流畅
	VideoQuality480p  = 2 // 640x480，标清
	VideoQuality576p  = 3 // 720x576，高清
	VideoQuality720p  = 4 // 1280x720, 超高清
	VideoQuality1080p = 5 // 1920x1080，蓝光
)

const (
	MediaStatusOnLine  = true // 上线
	MediaStatusNotOnLine = false // 下线
)

const (
	GenderUnknown = 0 // 未知状态
	GenderMale    = 1 // 男
	GenderFemale  = 2 // 女
)

const (
	CpTypeUnknown = 0        // 未知cp类型
	CpTypeVideo   = (1 << 0) // 视频cp
	CpTypeStream  = (1 << 1) //直播cp
)

const (
	PlayStatStart = "play_start"
	PlayStatStop  = "play_stop"
	PlayStatError = "play_error"
)

// redis配置类型
const (
	RedisTypeCommon = 1 // 通用缓存
	RedisTypeAd     = 3 // 广告统计
)

// 播放错误上报类型
const (
	PlayExceptionError     = 1 // 播放错误
	PlayExceptionDelay     = 2 // 起播超时次数过多
	PlayExceptionBuffering = 3 // 卡顿次数过多
)
