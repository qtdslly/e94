package main

import (
	"net/http"
	"background/common/logger"
	"background/others/config"
	"background/common/constant"
	"background/common/util"
	"io/ioutil"
	"strings"
	"errors"
	"github.com/tidwall/gjson"
	"background/others/model"
	"flag"
	"golang.org/x/text/encoding/simplifiedchinese"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"os"
	"bufio"
	"io"
	"time"
)

func main(){
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

	//fileName := "/root/Git/e94/src/background/others/tools/words.txt"
	//GetDominByFileWord(fileName,db)
	GetDominByRule(db)
}

func GetDominByRule(db *gorm.DB){
	zuis := []string{".com",".cn",".net",".cc"}
	charTypes := []string{"A","AA","AAA","AAAA","AAAAA","AAAAAA","AAAAAAA","AB","ABC","ABCD","AAAAB","AAABB","AABBB","ABBBB","ABABA","ABABB","ABAAA","ABBBB","AAAAB","AAABBB","ABABAB"}

	for _,zui := range zuis{
		for _,charType := range charTypes{
			GetDominByCharType(charType,zui,db)
		}
	}
}

func GetDominByCharType(charType string,zui string,db *gorm.DB){
	logger.Debug(charType,zui)
	chars := "a b c d e f g h i j k l m n o p q r s t u v w x y z"
	cs := strings.Split(chars," ")
	if charType == "A"{
		for _, ss := range cs{
			url := ss + zui
			getBaiDuDomin("",url,db)
		}
	}else if charType == "AA"{
		for _,ss := range cs{
			url := ss + ss + zui
			getBaiDuDomin("",url,db)
		}
	}else if charType == "AAA"{
		for _,ss := range cs{
			url := ss + ss + ss + zui
			getBaiDuDomin("",url,db)
		}
	}else if charType == "AAAA"{
		for _,ss := range cs{
			url := ss + ss + ss + ss + zui
			getBaiDuDomin("",url,db)
		}
	}else if charType == "AAAAA"{
		for _,ss := range cs{
			url := ss + ss + ss + ss + ss + zui
			getBaiDuDomin("",url,db)
		}
	}else if charType == "AAAAAA"{
		for _,ss := range cs{
			url := ss + ss + ss + ss + ss + ss + zui
			getBaiDuDomin("",url,db)
		}
	}else if charType == "AAAAAAA"{
		for _,ss := range cs{
			url := ss + ss + ss + ss + ss + ss + ss + zui
			getBaiDuDomin("",url,db)
		}
	}else if charType == "AB"{
		for _,ss := range cs{
			for _,mm := range cs{
				url := ss + mm + zui
				getBaiDuDomin("",url,db)
			}
		}
	}else if charType == "ABC"{
		for _,ss := range cs{
			for _,mm := range cs{
				for _,nn := range cs{
					url := ss + mm + nn + zui
					getBaiDuDomin("",url,db)
				}
			}
		}
	}else if charType == "ABCD"{
		for _,ss := range cs{
			for _,mm := range cs{
				for _,nn := range cs{
					for _, kk := range cs{
						url := ss + mm + nn + kk + zui
						getBaiDuDomin("",url,db)
					}
				}
			}
		}
	}else if charType == "AAAAB"{
		for _,ss := range cs{
			for _,mm := range cs{
				url := ss + ss + ss + ss + mm + zui
				getBaiDuDomin("",url,db)
			}
		}
	}else if charType == "AAABB"{
		for _,ss := range cs{
			for _,mm := range cs{
				url := ss + ss + ss + mm + mm + zui
				getBaiDuDomin("",url,db)
			}
		}
	}else if charType == "AABBB"{
		for _,ss := range cs{
			for _,mm := range cs{
				url := ss + ss + mm + mm + mm + zui
				getBaiDuDomin("",url,db)
			}
		}
	}else if charType == "ABBBB"{
		for _,ss := range cs{
			for _,mm := range cs{
				url := ss + mm + mm + mm + mm + zui
				getBaiDuDomin("",url,db)
			}
		}
	}else if charType == "ABABA"{
		for _,ss := range cs{
			for _,mm := range cs{
				url := ss + mm + ss + mm + ss + zui
				getBaiDuDomin("",url,db)
			}
		}
	}else if charType == "ABABB"{
		for _,ss := range cs{
			for _,mm := range cs{
				url := ss + mm + ss + mm + mm + zui
				getBaiDuDomin("",url,db)
			}
		}
	}else if charType == "ABAAA"{
		for _,ss := range cs{
			for _,mm := range cs{
				url := ss + mm + ss + ss + ss + zui
				getBaiDuDomin("",url,db)
			}
		}
	}else if charType == "ABBBB"{
		for _,ss := range cs{
			for _,mm := range cs{
				url := ss + mm + mm + mm + mm + zui
				getBaiDuDomin("",url,db)
			}
		}
	}else if charType == "AAAAB"{
		for _,ss := range cs{
			for _,mm := range cs{
				url := ss + ss + ss + ss + mm + zui
				getBaiDuDomin("",url,db)
			}
		}
	}else if charType == "AAABBB"{
		for _,ss := range cs{
			for _,mm := range cs{
				url := ss + ss + ss + mm + mm + mm + zui
				getBaiDuDomin("",url,db)
			}
		}
	}else if charType == "ABABAB"{
		for _,ss := range cs{
			for _,mm := range cs{
				url := ss + mm + ss + mm + ss + mm + zui
				getBaiDuDomin("",url,db)
			}
		}
	}
}

func GetDominByFileWord(fileName string,db *gorm.DB){
	f, err := os.Open(fileName)
	if err != nil {
		logger.Error(err)
		return
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		line1,err := DecodeToGBK(line)
		if err != nil{
			continue
		}
		line1 = strings.Replace(line1, "\n", "", -1)

		words := strings.Split(line1, "\t")
		for _ , word := range words{
			pin := util.TitleToPinyin(word)
			logger.Debug(word,pin)
			if !getDominByPinYin(word,pin,db){
				continue
			}

			fullpin := util.TitleToFullPinyin(word)
			logger.Debug(word,fullpin)

			if !getDominByPinYin(word,fullpin,db){
				continue
			}
		}
	}
}

func getDominByPinYin(word,pinyin string,db *gorm.DB)(bool){
	if !getBaiDuDomin(word,pinyin + ".com",db){return false}
	if !getBaiDuDomin(word,pinyin + ".cn",db){return false}
	if !getBaiDuDomin(word,pinyin + ".net",db){return false}
	if !getBaiDuDomin(word,pinyin + ".cc",db){return false}
	return true
}


func getBaiDuDomin(word,url string,db *gorm.DB)(bool){
	logger.Debug("=============================url:",url)
	var domain model.Domain
	domain.Url = url
	if err := db.Where("url = ?",domain.Url).First(&domain).Error ; err == nil{
		return true
	}
	status := getStatus(url)

	err,registerDate,expirationDate,registrarUrl,reseller,registrantCity,registrantProvince,email,phone,registrant,sponsoring,country,street := getDomainDetail(url)
	if err != nil{
		logger.Error(err)
		return false
	}

	logger.Debug(registerDate)
	logger.Debug(expirationDate)
	logger.Debug(registrarUrl)
	logger.Debug(reseller)
	logger.Debug(registrantCity)
	logger.Debug(registrantProvince)
	logger.Debug(email)
	logger.Debug(phone)
	logger.Debug(registrant)
	logger.Debug(sponsoring)
	logger.Debug(country)
	logger.Debug(street)

	domain.Status = uint32(status)
	domain.ExpirationDate = expirationDate
	domain.RegisterDate = registerDate
	domain.RegistrantCity = registrantCity
	domain.RegistrantProvince = registrantProvince
	domain.RegistrarUrl = registrarUrl
	domain.Reseller = reseller
	domain.RegistrantEmail = email
	domain.RegistrantPhone = phone
	domain.Chinese = word
	domain.Sponsoring = sponsoring
	domain.RegistrantName = registrant
	domain.RegistrantCountry = country
	domain.RegistrantStreet = street
	logger.Debug(domain.RegistrantCity)
	if err := db.Save(&domain).Error ; err != nil{
		logger.Error(err)
		return false
	}
	return true
}
func getStatus(url string)(int){
	time.Sleep(time.Second * 15)

	recv := getBaiDuDomainApiInfo(1,url)
	if recv == ""{
		return constant.DomainStatusUnknown
	}

	success := gjson.Get(recv, "success").Bool()
	status := gjson.Get(recv, "status").Int()

	if !success && status != 200{
		return constant.DomainStatusUnknown
	}

	accurate := gjson.Get(recv, "result.accurate")

	var domainStatus string
	if accurate.Exists() {
		re := accurate.Array()
		for _, v := range re {
			domainStatus = v.Get("status").String()
			break
		}
	}

	if domainStatus == "REGISTERED"{
		return constant.DomainStatusRegistered
	}else if domainStatus == "UNREGISTERED"{
		return constant.DomainStatusUnRegistered
	}

	return constant.DomainStatusUnknown
}

func getDomainDetail(url string)(error,string,string,string,string,string,string,string,string,string,string,string,string){
	time.Sleep(time.Second * 15)
	recv := getBaiDuDomainApiInfo(2,url)
	if recv == ""{
		return errors.New("调用百度云接口失败!"),"","","","","","","","","","","",""
	}

	recv = strings.Replace(recv,"\\r","",-1)
	success := gjson.Get(recv, "success").Bool()
	status := gjson.Get(recv, "status").Int()

	if !success && status != 200{
		return errors.New("百度云返回失败!"),"","","","","","","","","","","",""
	}

	registerDate := gjson.Get(recv, "result.data.registrationDate").String()
	expirationDate := gjson.Get(recv, "result.data.expirationDate").String()
	rawData := gjson.Get(recv, "result.data.rawData")

	var registrarUrl,reseller,registrantCity,registrantProvince,email,phone,registrant,sponsoring,country,street string
	if rawData.Exists() {
		re := rawData.Array()
		for _, v := range re {
			logger.Debug(v.String())
			if strings.Contains(v.String(),"Registrar URL"){
				registrarUrl = v.String()
				registrarUrl = strings.Replace(registrarUrl,"Registrar URL:","",-1)
			}

			if strings.Contains(v.String(),"Reseller"){
				reseller = v.String()
				reseller = strings.Replace(reseller,"Reseller:","",-1)
			}

			if strings.Contains(v.String(),"Registrant Name"){
				registrant = v.String()
				registrant = strings.Replace(registrant,"Registrant Name:","",-1)
			}

			if strings.Contains(v.String(),"Sponsoring Registrar"){
				sponsoring = v.String()
				sponsoring = strings.Replace(sponsoring,"Sponsoring Registrar:","",-1)
			}

			if strings.Contains(v.String(),"Registrant Country"){
				country = v.String()
				country = strings.Replace(country,"Registrant Country:","",-1)
			}

			if strings.Contains(v.String(),"Registrant Street"){
				street = v.String()
				street = strings.Replace(street,"Registrant Street:","",-1)
			}

			if strings.Contains(v.String(),"Registrant City"){
				registrantCity = v.String()
				registrantCity = strings.Replace(registrantCity,"Registrant City:","",-1)
			}

			if strings.Contains(v.String(),"Registrant State/Province"){
				registrantProvince = v.String()
				registrantProvince = strings.Replace(registrantProvince,"Registrant State/Province:","",-1)
			}

			if strings.Contains(v.String(),"Registrant Email"){
				email = v.String()
				email = strings.Replace(email,"Registrant Email:","",-1)
			}

			if strings.Contains(v.String(),"Registrant Phone:"){
				phone = v.String()
				phone = strings.Replace(phone,"Registrant Phone:","",-1)
			}
		}
	}

	registerDate = strings.Replace(registerDate,"年","-",-1)
	registerDate = strings.Replace(registerDate,"月","-",-1)
	registerDate = strings.Replace(registerDate,"日","",-1)

	expirationDate = strings.Replace(expirationDate,"年","-",-1)
	expirationDate = strings.Replace(expirationDate,"月","-",-1)
	expirationDate = strings.Replace(expirationDate,"日","",-1)

	registrarUrl = strings.Trim(registrarUrl," ")
	reseller = strings.Trim(reseller," ")
	registrantCity = strings.Trim(registrantCity," ")
	registrantProvince = strings.Trim(registrantProvince," ")
	email = strings.Trim(email," ")
	phone = strings.Trim(phone," ")
	registrant = strings.Trim(registrant," ")
	country = strings.Trim(country," ")
	street = strings.Trim(street," ")

	return nil,registerDate,expirationDate,registrarUrl,reseller,registrantCity,registrantProvince,email,phone,registrant,sponsoring,country,street
}
func getBaiDuDomainApiInfo(apiType int,url string)(string){
	var apiUrl string
	var postString string
	time.Sleep(time.Second * 15)
	values := strings.Split(url,".")
	dom := values[0]
	zui := values[1]

	if apiType == 1{//状态
		apiUrl = "https://cloud.baidu.com/api/bcd/search/status"
		postString = "{\"domainNames\":[{\"label\":\"" + dom + "\",\"tld\":\"" + zui + "\"}]}"
	}else{//详情
		apiUrl = "https://cloud.baidu.com/api/bcd/whois/detail"
		postString = "{\"domain\":\"" + url + "\",\"type\":\"NORMAL\"}"
	}

	requ, err := http.NewRequest("POST", apiUrl, strings.NewReader(postString))
	requ.Header.Add("Host", "cloud.baidu.com")
	//requ.Header.Add("Referer", "https://cloud.baidu.com/product/bcd/search.html?keyword=ezhantao")
	requ.Header.Add("Host", "cloud.baidu.com")
	requ.Header.Add("Content-Type", "application/json;charset=UTF-8")
	requ.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3278.0 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(requ)
	if err != nil {
		logger.Error(err)
		return ""
	}

	defer resp.Body.Close()

	recv, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return ""
	}

	logger.Debug(string(recv))

	if strings.Contains(string(recv),"查询过于频繁，请稍后再试"){
		time.Sleep(time.Minute * 15)
	}
	return string(recv)
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

/*
status info

{
    "success": true,
    "status": 200,
    "result": {
        "accurate": [
            {
                "domainName": "ezhantao111111.com",
                "status": "UNREGISTERED",
                "displayLevel": "ACCURATE"
            }
        ],
        "common": [],
        "recommend": [],
        "others": []
    }
}
*/
/*
detailinfo


{
    "success": true,
    "status": 200,
    "result": {
        "code": 0,
        "data": {
            "domain": "ezhantao.com",
            "queryTime": "2018-06-20T15:13:43Z",
            "privacy": false,
            "sponsoringRegistrar": "HiChina Zhicheng Technology Ltd.",
            "registrationDate": "2014年07月10日",
            "expirationDate": "2020年07月10日",
            "domainStatus": [
                "ok"
            ],
            "nameServer": [
                "DNS13.HICHINA.COM",
                "DNS14.HICHINA.COM"
            ],
            "rawData": [
                "Domain Name: ezhantao.com\r",
                "Registry Domain ID: 1866292912_DOMAIN_COM-VRSN\r",
                "Registrar WHOIS Server: grs-whois.hichina.com\r",
                "Registrar URL: http://whois.aliyun.com\r",
                "Updated Date: 2018-06-10T16:21:41Z\r",
                "Creation Date: 2014-07-10T10:07:01Z\r",
                "Registrar Registration Expiration Date: 2020-07-10T10:07:01Z\r",
                "Registrar: HiChina Zhicheng Technology Ltd.\r",
                "Registrar IANA ID: 420\r",
                "Reseller:\r",
                "Domain Status: ok https://icann.org/epp#ok\r",
                "Registrant City: wu han shi\r",
                "Registrant State/Province: hu bei\r",
                "Registry Registrant ID: Not Available From Registry\r",
                "Name Server: DNS13.HICHINA.COM\r",
                "Name Server: DNS14.HICHINA.COM\r",
                "DNSSEC: unsigned\r",
                "Registrar Abuse Contact Email: DomainAbuse@service.aliyun.com\r",
                "Registrar Abuse Contact Phone: +86.95187\r",
                "URL of the ICANN WHOIS Data Problem Reporting System: http://wdprs.internic.net/\r",
                ">>>Last update of WHOIS database: 2018-06-20T15:13:43Z <<<\r",
                "\r",
                "For more information on Whois status codes, please visit https://icann.org/epp\r",
                "\r",
                "Important Reminder: Per ICANN 2013RAA`s request, Hichina has modified domain names`whois format of dot com/net/cc/tv, you could refer to section 1.4 posted by ICANN on http://www.icann.org/en/resources/registrars/raa/approved-with-specs-27jun13-en.htm#whois The data in this whois database is provided to you for information purposes only, that is, to assist you in obtaining information about or related to a domain name registration record. We make this information available \"as is,\" and do not guarantee its accuracy. By submitting a whois query, you agree that you will use this data only for lawful purposes and that, under no circumstances will you use this data to: (1)enable high volume, automated, electronic processes that stress or load this whois database system providing you this information; or (2) allow, enable, or otherwise support the transmission of mass unsolicited, commercial advertising or solicitations via direct mail, electronic mail, or by telephone.  The compilation, repackaging, dissemination or other use of this data is expressly prohibited without prior written consent from us. We reserve the right to modify these terms at any time. By submitting this query, you agree to abide by these terms.For complete domain details go to:http://whois.aliyun.com/whois/domain/hichina.com",
                "",
                "",
                "   Domain Name: EZHANTAO.COM\r",
                "   Registry Domain ID: 1866292912_DOMAIN_COM-VRSN\r",
                "   Registrar WHOIS Server: grs-whois.hichina.com\r",
                "   Registrar URL: http://www.net.cn\r",
                "   Updated Date: 2018-06-10T16:21:41Z\r",
                "   Creation Date: 2014-07-10T10:07:01Z\r",
                "   Registry Expiry Date: 2020-07-10T10:07:01Z\r",
                "   Registrar: HiChina Zhicheng Technology Ltd.\r",
                "   Registrar IANA ID: 420\r",
                "   Registrar Abuse Contact Email: DomainAbuse@service.aliyun.com\r",
                "   Registrar Abuse Contact Phone: +86.95187\r",
                "   Domain Status: ok https://icann.org/epp#ok\r",
                "   Name Server: DNS13.HICHINA.COM\r",
                "   Name Server: DNS14.HICHINA.COM\r",
                "   DNSSEC: unsigned\r",
                "   URL of the ICANN Whois Inaccuracy Complaint Form: https://www.icann.org/wicf/\r",
                ">>> Last update of whois database: 2018-06-20T15:13:32Z <<<\r",
                "\r",
                "For more information on Whois status codes, please visit https://icann.org/epp\r",
                "\r",
                "NOTICE: The expiration date displayed in this record is the date the\r",
                "registrar's sponsorship of the domain name registration in the registry is\r",
                "currently set to expire. This date does not necessarily reflect the expiration\r",
                "date of the domain name registrant's agreement with the sponsoring\r",
                "registrar.  Users may consult the sponsoring registrar's Whois database to\r",
                "view the registrar's reported date of expiration for this registration.\r",
                "\r",
                "TERMS OF USE: You are not authorized to access or query our Whois\r",
                "database through the use of electronic processes that are high-volume and\r",
                "automated except as reasonably necessary to register domain names or\r",
                "modify existing registrations; the Data in VeriSign Global Registry\r",
                "Services' (\"VeriSign\") Whois database is provided by VeriSign for\r",
                "information purposes only, and to assist persons in obtaining information\r",
                "about or related to a domain name registration record. VeriSign does not\r",
                "guarantee its accuracy. By submitting a Whois query, you agree to abide\r",
                "by the following terms of use: You agree that you may use this Data only\r",
                "for lawful purposes and that under no circumstances will you use this Data\r",
                "to: (1) allow, enable, or otherwise support the transmission of mass\r",
                "unsolicited, commercial advertising or solicitations via e-mail, telephone,\r",
                "or facsimile; or (2) enable high volume, automated, electronic processes\r",
                "that apply to VeriSign (or its computer systems). The compilation,\r",
                "repackaging, dissemination or other use of this Data is expressly\r",
                "prohibited without the prior written consent of VeriSign. You agree not to\r",
                "use electronic processes that are automated and high-volume to access or\r",
                "query the Whois database except as reasonably necessary to register\r",
                "domain names or modify existing registrations. VeriSign reserves the right\r",
                "to restrict your access to the Whois database in its sole discretion to ensure\r",
                "operational stability.  VeriSign may restrict or terminate your access to the\r",
                "Whois database for failure to abide by these terms of use. VeriSign\r",
                "reserves the right to modify these terms at any time.\r",
                "\r",
                "The Registry database contains ONLY .COM, .NET, .EDU domains and\r",
                "Registrars.\r"
            ]
        }
    }
}

*/