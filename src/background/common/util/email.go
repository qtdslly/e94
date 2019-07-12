package util

import (
	"background/common/logger"
	"github.com/go-gomail/gomail"
)

func SendEmail(subject,content string,toUser string)(error){
	m := gomail.NewMessage()
	m.SetAddressHeader("From", "202037514@qq.com", "e站淘")  // 发件人
	m.SetHeader("To",  // 收件人
		m.FormatAddress(toUser, "e站淘"),
	)
	m.SetHeader("Subject", subject)  // 主题
	m.SetBody("text/html", content)  // 正文

	//d := gomail.NewPlainDialer("smtp.aliyun.com", 465, "qtdslly@aliyun.com", "qtds@lly413")  // 发送邮件服务器、端口、发件人账号、发件人密码
  //d := gomail.NewPlainDialer("smtp.aliyun.com", 465, "ezhantao@aliyun.com", "hacker@LLY413")  // 发送邮件服务器、端口、发件人账号、发件人密码
  d := gomail.NewPlainDialer("smtp.qq.com", 465, "202037514@qq.com", "yquusrfaezctbibj")  // 发送邮件服务器、端口、发件人账号、发件人密码

  if err := d.DialAndSend(m); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
