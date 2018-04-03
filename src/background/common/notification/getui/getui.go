package getui

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 常量
const (
	PUSH_NUM_KEY                                                   = "med-msg-svr:push-number:"
	GE_TUI_CONF_PATH                                        string = "../config/getui.yaml"
	IS_OFFLINE_ANDROID, IS_OFFLINE_IOS                      bool   = true, true
	PHONE_TYPE_ANDROID, PHONE_TYPE_IOS                      string = "ANDROID", "IOS"
	PUSH_NETWORK_TYPE_ANDROID, PUSH_NETWORK_TYPE_IOS        byte   = 0, 0
	PUSH_MESSAGE_SPEED_ANDROID, PUSH_MESSAGE_SPEED_IOS      uint32 = 3000, 3000
	PUSH_SINGLE, PUSH_GROUP_LIMIT_MIN, PUSH_GROUP_LIMIT_MAX int    = 1, 2, 10
	OFFLINE_EXPIRE_TIME_ANDROID, OFFLINE_EXPIRE_TIME_IOS    uint32 = 300 * 12 * 1000, 300 * 12 * 1000
)

// 个推配置
type GeTuiConf struct {
	Host         string `json:"host"`
	AppId        string `json:"app_id"`
	AppKey       string `json:"app_key"`
	MasterSecret string `json:"master_secret"`
}

type Message struct {
	AppKey            string `json:"appkey"`
	IsOffline         bool   `json:"is_offline, omitempty"`          //是否离线推送
	OfflineExpireTime uint32 `json:"offline_expire_time, omitempty"` //消息离线存储有效期，单位：ms
	PushNetworkType   uint32 `json:"push_network_type, omitempty"`   //选择推送消息使用网络类型，0：不限制，1：wifi
	MsgType           string `json:"msgtype"`                        //消息应用类型，可选项：notification、link、notypopload、transmission
}

type Transmission struct {
	TransmissionType    bool   `json:"transmission_type, omitempty"` //收到消息是否立即启动应用，true为立即启动，false则广播等待启动，默认是否
	TransmissionContent string `json:"transmission_content"`         //透传内容
	//DurationBegin       string    `json:"duration_begin, omitempty"`    //设定展示开始时间，格式为yyyy-MM-dd HH:mm:ss
	//DurationEnd         string    `json:"duration_end, omitempty"`      //设定展示结束时间，格式为yyyy-MM-dd HH:mm:ss
}

//透传消息內容
type Payload struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Cover   string `json:"cover"`
	Url     string `json:"url"`
}

type PushInfo struct {
	Aps        Aps      `json:"aps"`
	Multimedia []*Media `json:"multimedia"`
}

type Aps struct {
	Alert            Alert  `json:"alert"`             //消息
	AutoBadge        string `json:"autoBadge"`         //用于计算icon上显示的数字，还可以实现显示数字的自动增减，如“+1”、 “-1”、 “1” 等，计算结果将覆盖badge
	Sound            string `json:"sound"`             //通知铃声文件名，无声设置为“com.gexin.ios.silence”
	ContentAvailable uint32 `json:"content-available"` //推送直接带有透传数据
	Category         string `json:"category"`          //在客户端通知栏触发特定的action和button显示
}

type Alert struct {
	Body            string   `json:"body"`              //通知文本消息
	ActionLocKey    string   `json:"action-loc-key"`    //（用于多语言支持）指定执行按钮所使用的Localizable.strings
	LocKey          string   `json:"loc-key"`           //（用于多语言支持）指定Localizable.strings文件中相应的key
	LocArgs         []string `json:"loc-args"`          //如果loc-key中使用了占位符，则在loc-args中指定各参数
	LaunchImage     string   `json:"launch-image"`      //指定启动界面图片名
	Title           string   `json:"title"`             //通知标题
	TitleLocKey     string   `json:"titile-loc-key"`    //(用于多语言支持）对于标题指定执行按钮所使用的Localizable.strings,仅支持iOS8.2以上版本
	TitleLocArgs    []string `json:"title-loc-args"`    //对于标题,如果loc-key中使用的占位符，则在loc-args中指定各参数,仅支持iOS8.2以上版本
	SubTitle        string   `json:"subtitle"`          //通知子标题,仅支持iOS8.2以上版本
	SubTitleLocKey  string   `json:"subtitle-loc-key"`  //当前本地化文件中的子标题字符串的关键字,仅支持iOS8.2以上版本
	SubTitleLocArgs []string `json:"subtitle-loc-args"` //当前本地化子标题内容中需要置换的变量参数 ,仅支持iOS8.2以上版本
}

type Media struct {
	Url      string `json:"url"`       //必须，多媒体资源地址
	Type     uint32 `json:"type"`      //必须，资源类型（1.图片，2.音频， 3.视频）
	OnlyWifi bool   `json:"only_wifi"` //是否只在wifi环境下加载，如果设置成true,但未使用wifi时，会展示成普通通知
}

type Style struct {
	Type        int    `json:"type"`         //必须，固定为0
	Text        string `json:"text"`         //必须，通知内容
	Title       string `json:"title"`        //必须，通知标题
	Logo        string `json:"logo"`         //必须，通知的图标名称，包含后缀名（需要在客户端开发时嵌入），如“push.png”
	IsRing      bool   `json:"is_ring"`      //收到通知是否响铃：true响铃，false不响铃。默认响铃
	IsVibrate   bool   `json:"is_vibrate"`   //收到通知是否振动：true振动，false不振动。默认振动
	IsClearable bool   `json:"is_clearable"` //通知是否可清除： true可清除，false不可清除。默认可清除
}

type Condition struct {
	Key     string   `json:"key"`      //筛选条件类型名称(省市region,手机类型phonetype,用户标签tag)
	Values  []string `json:"values"`   //筛选参数
	OptType uint32   `json:"opt_type"` //筛选参数的组合，0:取参数并集or，1：交集and，2：相当与not in {参数1，参数2，....}
}

// 取得个推配置
var GetGeTuiConf = func(appId, appKey, masterSecret string) *GeTuiConf {
	geTuiConf := GeTuiConf{}
	geTuiConf.Host = "http://sdk.open.api.igexin.com/apiex.htm"
	geTuiConf.AppId = appId
	geTuiConf.AppKey = appKey
	geTuiConf.MasterSecret = masterSecret
	return &geTuiConf
}

func GetCurrentTime() int64 {
	t := time.Now().Unix() * 1000
	return t
}

func HttpPostJson(conf *GeTuiConf, url string, params map[string]interface{}, headers map[string]string) map[string]interface{} {
	ret := httpPost(url, params, headers)
	if ret["result"] == "sign_error" {
		connect(conf)
		ret = httpPost(url, params, headers)
	}

	return ret
}

func httpPost(url string, params map[string]interface{}, headers map[string]string) map[string]interface{} {
	data, _ := json.Marshal(params)
	fmt.Printf("-------params---------%s\n", data)
	payload := strings.NewReader(string(data))

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/json")
	for k, v := range headers {
		if k == "authtoken" {
			req.Header["authtoken"] = []string{v}
		} else {
			req.Header.Add(k, v)
		}
	}

	tryTime := 1
tryAgain:
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("第"+strconv.Itoa(tryTime)+"次", "请求失败")
		tryTime += 1
		if tryTime < 4 {
			goto tryAgain
		}
		return map[string]interface{}{"result": "post error"}
	}
	// close body read before return
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	var ret map[string]interface{}
	json.Unmarshal(body, &ret)
	return ret
}

func connect(conf *GeTuiConf) bool {
	sign := getSign(conf.AppKey, GetCurrentTime(), conf.MasterSecret)
	params := map[string]interface{}{}
	params["action"] = "connect"
	params["appkey"] = conf.AppKey
	params["timeStamp"] = GetCurrentTime()
	params["sign"] = sign

	rep := httpPost(conf.Host, params, map[string]string{})
	fmt.Println("rep")
	fmt.Println(rep)
	if "success" == rep["result"] {
		return true
	} else {
		fmt.Println("connect failed")
		return false
	}
}

func getSign(appKey string, timeStamp int64, masterSecret string) string {
	rawValue := appKey + strconv.FormatInt(timeStamp, 10) + masterSecret
	h := md5.New()
	io.WriteString(h, rawValue)
	return hex.EncodeToString(h.Sum(nil))
}
