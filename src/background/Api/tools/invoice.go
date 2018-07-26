package main

import (
	"net/http"
	"background/common/logger"
	"io/ioutil"
	_ "github.com/go-sql-driver/mysql"
	"background/Api/config"
)

type TaxpayerInfo struct {
	Code	string `json:"code"`        //发票代码
	Number	string `json:"number"`      //发票号码
	Type	string `json:"type"`        //发票类型
	RecNum  string `json:"rec_num"`     //纳税人识别号
	Company	string `json:"company"`     //开票单位
}
func main(){
	logger.SetLevel(config.GetLoggerLevel())

	var tax TaxpayerInfo
	tax.Code = "121001121071"
	tax.Number = "00722838"
	tax.Type = "0"
	tax.RecNum = "210114573484030"
	tax.Company = "沈阳业乔新业汽车销售服务有限公司"
}

func GetInvoiceInfo(tax TaxpayerInfo){
	apiUrl := "http://fapiao.youshang.com/seekUrlAction.do?querytype=doseek&cityid=12100&invoiceTypeId=1210001&fpdm=121001121071&fphm=00722838&fplx=0&nsrsbh=210114573484030&nsrmc=%E6%B2%88%E9%98%B3%E4%B8%9A%E4%B9%94%E6%96%B0%E4%B8%9A%E6%B1%BD%E8%BD%A6%E9%94%80%E5%94%AE%E6%9C%8D%E5%8A%A1%E6%9C%89%E9%99%90%E5%85%AC%E5%8F%B8"
	requ, err := http.NewRequest("GET", apiUrl, nil)
	requ.Header.Add("app_version", "1.0.0")
	//requ.Header.Add("Referer", "https://cloud.baidu.com/product/bcd/search.html?keyword=ezhantao")
	requ.Header.Add("app_key", "5Hwf3G51ZUPT")
	requ.Header.Add("installation_id","1807150014158383")
	requ.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3278.0 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(requ)
	if err != nil {
		logger.Error(err)
		return
	}

	defer resp.Body.Close()

	recv, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(string(recv))
}