package main

import (
  "background/ezhantao/config"
  "background/common/logger"
  "background/common/util"
  "time"
  "os"
  "bufio"
  "io"
  "strings"
)

func main(){
  logger.SetLevel(config.GetLoggerLevel())

  subject := "免费领福利啦"
  content := "<img src='http://www.ezhantao.com/html/images/jc.jpg'/>"
  //content := "<div style='text-indent:2em;'><h3>关注公众号e站淘，在手淘选中宝贝名称并复制到此公众号，然后发送宝贝名称，系统返回优惠券的淘口令，复制口令到手机淘宝即可领取优惠券，支付时选择此优惠券即可享受相应的折扣</h3></div>" +
  //  "<div><h2 style='text-align:left;color:red;'>第一步 微信搜索公众号ezhantao-com并关注</h2></br><img src='http://www.ezhantao.com/html/images/zero.jpg'/></div></br></br></br>" +
  //  "<div><h2 style='text-align:left;color:red;'>第二步 打开手淘复制宝贝名称</h2></br><img src='http://www.ezhantao.com/html/images/first.png'/></div></br></br></br>" +
  //  "<div><h2 style='text-align:left;color:red;'>第三步 发送宝贝名称到公众号获取淘口令</h2></br><img src='http://www.ezhantao.com/html/images/second.png'/></div></br></br></br>" +
  //  "<div><h2 style='text-align:left;color:red;'>第四步 复制淘口令到手淘点击打开按钮</h2></br><img src='http://www.ezhantao.com/html/images/third.png'/></div></br></br></br>" +
  //  "<div><h2 style='text-align:left;color:red;'>第五步 在粉丝福利购页面领取优惠券</h2></br><img src='http://www.ezhantao.com/html/images/four.png'/></div>"

  //addrs := []string{"86282864","592189599","2680656796","2548257202","415382006"}
  //for _ ,addr := range addrs{
  //  util.SendEmail(subject,content,addr)
  //  time.Sleep(time.Second * 3)
  //}



  var err error
  //f, err := os.Open("/home/lyric/Git/e94/src/background/newmovie/tools/stream.txt")
  //f, err := os.Open("/root/Git/e94/src/background/newmovie/tools/stream.txt")
  f, err := os.Open("C:/work/code/e94/src/background/ezhantao/tools/QQList.txt")

  if err != nil {
    logger.Error(err)
    return
  }
  defer f.Close()

  rd := bufio.NewReader(f)
  for {
    line, err := rd.ReadString('\n')
    if err != nil || io.EOF == err {
      logger.Error(err)
      break
    }
    addr := strings.Replace(line, "\n", "", -1)
    if addr == "86282864" || addr == "592189599" || addr == "2680656796" || addr == "2548257202" || addr == "415382006" {
      continue
    }
    addr = addr + "@qq.com"
    //logger.Debug(addr)

    err = util.SendEmail(subject,content,addr)
    if err == nil{
      logger.Debug("发送邮件给" + addr + "成功")
    }else{
      if !strings.Contains(err.Error() ,"access denied"){
        time.Sleep(time.Hour)
      }
    }
    return
    time.Sleep(time.Second * 15)
  }

}
