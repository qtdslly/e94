package common

import (
	"background/common/logger"
	"github.com/go-gomail/gomail"
)

func SendEmail(subject,content string)(bool){
	m := gomail.NewMessage()
	m.SetAddressHeader("From", "qtdslly@aliyun.com", "Lyric")  // 发件人
	m.SetHeader("To",  // 收件人
		m.FormatAddress("947883972@qq.com", "Lyric"),
	)
	m.SetHeader("Subject", subject)  // 主题
	m.SetBody("text/html", content)  // 正文

	d := gomail.NewPlainDialer("smtp.aliyun.com", 465, "qtdslly@aliyun.com", "qtds@lly413")  // 发送邮件服务器、端口、发件人账号、发件人密码
	if err := d.DialAndSend(m); err != nil {
		logger.Error(err)
		return false
	}
	return true
}
