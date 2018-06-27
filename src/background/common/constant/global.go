package constant

const (
	MaxPageSize     = 256
	DefaultPageSize = 10

	SortAsc       = "asc"
	SortDesc      = "desc"
	DefaultColumn = "id"

	CdnSettingKey = "cdn_setting_key"

	// Kv store keys
	ResourceGroupSyncedAt   = "resource_group_synced_at"
	IptvVideoGroupSyncedAt  = "iptv_video_group_synced_at"
	IptvSeriesGroupSyncedAt = "iptv_series_group_synced_at"
	PersonSyncedAt          = "person_synced_at"
	PersonMediaSyncedAt     = "personmedia_synced_at"
	LanguageSyncedAt        = "language_synced_at"
	TranslationSyncedAt     = "translation_synced_at"
	DongFangMediaSyncedAt   = "dongfang_media_synced_at"
	ChinaCloudMediaSyncedAt = "chinacloud_media_synced_at"

	WhitelistSettingKey    = "whitelist_setting_key"
	CommentSettingKey      = "comment_setting_key"
	FilterSettingKey       = "filter_setting_key"
	ConfigSettingKey       = "config_setting_key"
	AppProviderSettingKey  = "app_provider_setting_key"
	LibrarySettingKey      = "library_setting_key"
	RedisSettingKey        = "redis_setting_key"
	SmsSettingKey          = "sms_setting_key"
	ThirdPartyInfoKey      = "third_party_key"
	PlayStatKey            = "play_stat_key"
	ChannelCacheKey        = "channel_cache_key"
	NotificationSettingKey = "notification_setting_key"
	FunctionSettingKey     = "function_setting_key"

	TideSyncTime         = "tide_sync_time"
	TideIncreaseSyncTime = "tide_increase_sync_time"

	TmpStorage      = "tmp"
	VideoTmpStorage = "video_tmp"
	LogStorage      = "log"

	DefaultLanguageCode = "zh"

	InternalTypeAll   = 0 //全部
	InternalTypeTrue  = 1 //内部商品
	InternalTypeFalse = 2 //非内部商品

	ActiveTypeAll   = 0 //全部
	ActiveTypeTrue  = 1 //已上架
	ActiveTypeFalse = 2 //已下架

	DisabledTypeAll   = 0 //全部
	DisabledTypeTrue  = 1 //禁用
	DisabledTypeFalse = 2 //可用

	OnlineTypeAll   = 0 //全部
	OnlineTypeTrue  = 1 //已上线
	OnlineTypeFalse = 2 //未上线

	DownloadableTypeAll   = 0 //全部
	DownloadableTypeTrue  = 1 //可下载
	DownloadableTypeFalse = 2 //不可下载

	SystemTypeAll   = 0 //全部分组
	SystemTypeTrue  = 1 //系统分组
	SystemTypeFalse = 2 //非系统分组

	OsTypeAll = 0 //全部系统类型

	UserCacheTTL     = 300     // 用户结构缓存时间，单位秒
	CommentCacheTTL  = 300     // 评论缓存时间
	AreaFeedCacheTTL = 60 * 30 // 地域消息流缓存时间

	SearchTypeUnknown    = 0 //未知类型
	SearchTypeProduct    = 1 //商品搜索, search_value为商品id
	SearchTypePackage    = 2 //套餐搜索, search_value为套餐id
	SearchTypeVipProduct = 3 //vip商品, search_value留空
	SearchTypeVipPackage = 4 //vip套餐, search_value留空
	SearchTypePrivilege  = 5 //特权商品搜索, search_value为商品的特权

	ReachedTypeAll   = 0
	ReachedTypeTrue  = 1
	ReachedTypeFalse = 2

	ExpiredTypeAll   = 0
	ExpiredTypeTrue  = 1
	ExpiredTypeFalse = 2

	MobileTVAppId = 1
	ChinaTVAppId  = 11

	AccessKey       = "AccessKey"
	ModuleSignature = "Signature"
	ModuleSalt      = "module_salt"

	PropertyArtist    = "艺术家"
	PropertyWriter    = "作者"
	PropetyCategoryId = 5

	SmsProviderYunPian    = 1
	SmsProviderYunTongXin = 2

	ScriptSettingKey = "script_setting_key"

	ContentProviderSystem   = 0
	ContentProviderYouKu    = 1
	ContentProviderMgtv     = 2
	ContentProviderIqiyi    = 3
	ContentProviderTencent  = 4
	ContentProviderSohu     = 5
	ContentProviderDouYin   = 6
	ContentProviderPear     = 7
	ContentProviderKuai     = 8
	ContentProviderMigu     = 9

	PictureProivderBiZhiJinXuan = 1

	DomainStatusUnknown      = 0
	DomainStatusRegistered   = 1
	DomainStatusUnRegistered = 2

	DouBanCrawlStatusSuccess = 0
	DouBanCrawlStatusError   = 1
	DouBanCrawlStatusReady   = 2

	DownloadUrlProviderDytt8 = 1

)
