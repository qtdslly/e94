package task


import (
  "net/http"
  "background/common/logger"
  "background/common/constant"
  "io/ioutil"
  "strings"
  "errors"
  "github.com/tidwall/gjson"
  "background/others/model"
  "golang.org/x/text/encoding/simplifiedchinese"

  _ "github.com/go-sql-driver/mysql"
  "github.com/jinzhu/gorm"
  "time"
  "strconv"

  "github.com/PuerkitoBio/goquery"

)

func GetUrl(db *gorm.DB){
  zuis := []string{".com",".cn",".cc"}
  chars := "a b c d e f g h i j k l m n o p q r s t u v w x y z"
  cs := strings.Split(chars," ")

  GetOneLength(zuis,cs,db)
  GetTwoLength(zuis,cs,db)
  GetThreeLength(zuis,cs,db)

  GetFourLength(zuis,cs,db)
  GetFiveLength(zuis,cs,db)
  GetSixLength(zuis,cs,db)
}

func GetOneLength(zuis,cs []string ,db *gorm.DB){
  var url string = ""
  for _,c := range cs{
    for _,zui := range zuis {
      url = c + zui
      CreateDomain(url, db)
    }
  }
}


func GetTwoLength(zuis,cs []string ,db *gorm.DB){
  var url string = ""
  for _,c := range cs{
    for _,c1 := range cs{
      for _,zui := range zuis {
        url = c + c1 + zui
        CreateDomain(url, db)
      }
    }
  }
}

func GetThreeLength(zuis,cs []string ,db *gorm.DB){
  var url string = ""
  for _,c := range cs{
    for _,c1 := range cs{
      for _,c2 := range cs{
        for _,zui := range zuis {

          url = c + c1 + c2 + zui
          CreateDomain(url, db)
        }
      }
    }
  }
}


func GetFourLength(zuis,cs []string ,db *gorm.DB){
  var url string = ""
  for _,c := range cs{
    for _,c1 := range cs{
      for _,c2 := range cs{
        for _,c3 := range cs{
          for _,zui := range zuis {

            url = c + c1 + c2 + c3 + zui
            CreateDomain(url, db)
          }
        }
      }
    }
  }
}

func GetFiveLength(zuis,cs []string ,db *gorm.DB){
  var url string = ""
  for _,c := range cs{
    for _,c1 := range cs{
      for _,c2 := range cs{
        for _,c3 := range cs{
          for _,c4 := range cs{
            for _,zui := range zuis {

              url = c + c1 + c2 + c3 + c4 + zui
              CreateDomain(url, db)
            }
          }
        }
      }
    }
  }
}

func GetSixLength(zuis,cs []string ,db *gorm.DB){
  var url string = ""
  for _,c := range cs{
    for _,c1 := range cs{
      for _,c2 := range cs{
        for _,c3 := range cs{
          for _,c4 := range cs{
            for _,c5 := range cs{
              for _,zui := range zuis {
                url = c + c1 + c2 + c3 + c4 + c5 + zui
                CreateDomain(url, db)
              }
            }
          }
        }
      }
    }
  }
}

func CreateDomain(url string,db *gorm.DB){
  var domain model.Domain
  domain.Url = url
  if err := db.Where("url = ?",url).First(&domain).Error ; err == nil{
    return
  }
  domain.IsGet = false
  if err := db.Create(&domain).Error ; err != nil{
    logger.Error(err)
    return
  }
}


func GetBaiDuDomin(db *gorm.DB){
  for{
    time.Sleep(time.Second * 20)
    var domain model.Domain
    if err := db.Where("status = 0 and is_get = 0").First(&domain).Error ; err != nil{
      logger.Error(err)
      return
    }

    status := getStatus(domain.Url)

    err,registerDate,expirationDate,registrarUrl,reseller,registrantCity,registrantProvince,email,phone,registrant,sponsoring,country,street := getDomainDetail(domain.Url)
    if err != nil{
      logger.Error(err)
      domain.IsGet = true
      if err := db.Save(&domain).Error ; err != nil{
        logger.Error(err)
        return
      }
      return
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
    domain.Chinese = ""
    domain.Sponsoring = sponsoring
    domain.RegistrantName = registrant
    domain.RegistrantCountry = country
    domain.RegistrantStreet = street
    domain.IsGet = true

    logger.Debug(domain.RegistrantCity)
    if err := db.Save(&domain).Error ; err != nil{
      logger.Error(err)
      return
    }
  }

}

func getStatus(url string)(int){
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
  recv := getBaiDuDomainApiInfo(2,url)
  if recv == ""{
    logger.Debug("调用百度云接口失败!")
    return errors.New("调用百度云接口失败!"),"","","","","","","","","","","",""
  }

  recv = strings.Replace(recv,"\\r","",-1)
  success := gjson.Get(recv, "success").Bool()
  status := gjson.Get(recv, "status").Int()

  if !success && status != 200{
    logger.Debug("调用百度云接口失败!")
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
    time.Sleep(time.Hour * 24)
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


func GetDominInfo(url string){
  apiurl := "https://whois.ename.net/" + url
  logger.Debug(apiurl)
  document, err := goquery.NewDocument(apiurl)
  if err != nil {
    logger.Error(err)

    return
  }

  block := document.Find(".nicehash").Eq(0).Find("tbody").Eq(0).Find("tr").Eq(2).Find("td").Eq(3).Text()

  block = strings.Split(block, " ")[0]
  logger.Debug(block)

  theory, err := strconv.Atoi(block)
  if err != nil {
    logger.Error(err)

    return
  }
  logger.Debug(theory)
}
