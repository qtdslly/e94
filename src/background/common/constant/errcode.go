package constant

import "strings"

const (
	UnknownError              = -1
	Success                   = 0
	Failure                   = 1
	StatusInternalServerError = 500

	// api side error, starts with 1000-
	UserNotExists          = 10001
	ProductNotExists       = 10002
	PackageNotExists       = 10003
	ProductSoldOut         = 10004
	PackageSoldOut         = 10005
	OrderProcessed         = 10006
	OrderCanceled          = 10007
	ProductIdentityMissing = 10008 // apple product identity is missing
	UnSupportedOsType      = 10009 // only android, ios are accepted
	InvalidPayMethod       = 10010
	TopicNotExists         = 10011
	SohuIdNotExists        = 10012
	RedeemCodeInvalid      = 10013
	RedeemCodeNotStarted   = 10014
	RedeemCodeExpired      = 10015
	RedeemCodeUsed         = 10016
	NoAvailableRedeemCode  = 10017
	UserLoginBanned        = 10018
	UserCommentBanned      = 10019
	ThirdPartyLoginFailed  = 10020
	SmsVerificationFailed  = 10021
	ThirdPartyUnBindError  = 10022
	ThirdPartyUnBindFailed = 10023
	ThirdPartyAlreadyBind  = 10024
	MobileAlreadyBind      = 10025
	CommentDisabled        = 10026
	AdConfigNotEnabled     = 10027
	RegistrationDisabled   = 10028
	UnAuthorizedResource   = 10029
	IncorrectBundleId      = 10030

	// ims side error, starts with 2000-
	AdminNotExists        = 20001
	ChangyanFailed        = 20002
	InvalidVcode          = 20003
	ApiNotRegisted        = 20004
	ModuleAccessDenied    = 20005
	InvalidAccessKey      = 20006
	InvalidAdminToken     = 20007
	AdminNotLogin         = 20008
	OperationNotPermitted = 20009
	ResourceGroupNotEmpty = 20010
	PropertyNotEmpty      = 20011
	TagNotEmpty           = 20012
	PageHasBeenRelated    = 20013
	PageHasChildren       = 20014
	PageGroupNotEmpty     = 20015
	UnAuthorizedGroupId   = 20016
	VipLevelRepeated      = 20017
	AppAccessDenied       = 20018
	InvalidInput          = 21000
	InvalidPlayurl        = 21001

	// common (both api + ims) error, starts with 3000-
	WrongUsernamePassword = 30002
)

func TranslateErrCode(code int, extra ...string) string {
	var msg string
	switch code {
	case UnknownError:
		msg = "Unknown error"
	case AdminNotExists:
		msg = "管理员不存在"
	case WrongUsernamePassword:
		msg = "用户名或密码错误"
	case AdminNotLogin:
		msg = "用户未登录"
	case InvalidVcode:
		msg = "验证码错误"
	case ChangyanFailed:
		msg = "畅言接口调用失败"
	case UserNotExists:
		msg = "用户不存在"
	case ProductNotExists:
		msg = "商品不存在"
	case PackageNotExists:
		msg = "套餐不存在"
	case ProductSoldOut:
		msg = "商品已售罄"
	case PackageSoldOut:
		msg = "套餐已售罄"
	case ApiNotRegisted:
		msg = "api对应的模块未找到"
	case ModuleAccessDenied:
		msg = "模块访问受限"
	case ProductIdentityMissing:
		msg = "商品没有苹果对应id"
	case UnSupportedOsType:
		msg = "不支持的系统类型,仅支持android/ios"
	case InvalidPayMethod:
		msg = "不支持的支付方式"
	case TopicNotExists:
		msg = "未找到对应Topic"
	case SohuIdNotExists:
		msg = "该用户的SohuId尚未设置，无法进行禁言相关操作，请在该用户下次发表评论后重试"
	case RedeemCodeInvalid:
		msg = "无效的激活码"
	case RedeemCodeNotStarted:
		msg = "激活码还没有开始"
	case RedeemCodeExpired:
		msg = "激活码已经过期"
	case RedeemCodeUsed:
		msg = "激活码已经被使用"
	case NoAvailableRedeemCode:
		msg = "没有可用的激活码"
	case UserLoginBanned:
		msg = "用户被禁止登陆"
	case UserCommentBanned:
		msg = "用户被禁止发言"
	case InvalidAccessKey:
		msg = "无效的访问秘钥"
	case InvalidAdminToken:
		msg = "无效的管理员令牌"
	case ThirdPartyLoginFailed:
		msg = "第三方账号登录失败"
	case ThirdPartyUnBindError:
		msg = "用户尚未绑定账号,不能解绑"
	case ThirdPartyUnBindFailed:
		msg = "用户只有一个绑定账号,不能再解绑"
	case SmsVerificationFailed:
		msg = "短信验证失败"
	case OperationNotPermitted:
		msg = "操作不被允许"
	case ThirdPartyAlreadyBind:
		msg = "第三方账号已绑定，无法二次绑定"
	case MobileAlreadyBind:
		msg = "手机号码已绑定，无法二次绑定"
	case ResourceGroupNotEmpty:
		msg = "资源组下面不为空"
	case PropertyNotEmpty:
		msg = "属性下面不为空"
	case TagNotEmpty:
		msg = "标签下面不为空"
	case InvalidInput:
		msg = "输入不合法"
	case CommentDisabled:
		msg = "评论未开启"
	case PageHasBeenRelated:
		msg = "页面已被关联使用"
	case PageHasChildren:
		msg = "页面包含子页面"
	case PageGroupNotEmpty:
		msg = "页面组不为空"
	case AdConfigNotEnabled:
		msg = "广告位配置没有开启"
	case UnAuthorizedGroupId:
		msg = "group_id未授权"
	case VipLevelRepeated:
		msg = "Vip level 重复"
	case AppAccessDenied:
		msg = "App访问受限"
	case RegistrationDisabled:
		msg = "用户注册功能被禁用"
	case UnAuthorizedResource:
		msg = "资源未授权"
	case IncorrectBundleId:
		msg = "Bundle Id不正确"
	case InvalidPlayurl:
		msg = "无效的播放链接"
	default:
	}

	if len(extra) > 0 {
		msg = msg + ": " + strings.Join(extra, ",")
	}
	return msg
}
