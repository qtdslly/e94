package main

import (
  "net/http"
  "io/ioutil"
  "background/wechart/config"
  "background/common/logger"

  "time"
  "fmt"
)

func main(){
  logger.SetLevel(config.GetLoggerLevel())
  p := time.Now().Unix()
  t := fmt.Sprintf("%d",p)
  apiUrl := "http://pub.alimama.com/common/code/getAuctionCode.json?auctionid=37016930337&adzoneid=adzoneid&siteid=584850166&scenes=1&t="+t+"&_tb_token_=qO2Nj1Sk4Rq&pvid=10_122.233.43.77_1118_1489238002348"
  requ, err := http.NewRequest("GET", apiUrl, nil)
  requ.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

  resp, err := http.DefaultClient.Do(requ)
  if err != nil {
    logger.Debug(err)
    return
  }

  recv, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    logger.Error(err)
    return
  }
  logger.Debug(string(recv))
}
