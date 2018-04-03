package main

import (
	"github.com/go-gomail/gomail"
)

func main() {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", "947883972@qq.com", "Lyric")  // 发件人
	m.SetHeader("To",  // 收件人
		m.FormatAddress("947883972@qq.com", "Lyric"),
	)
	m.SetHeader("Subject", "股票提示信息")  // 主题
	m.SetBody("text/html", "该买股票了 <a href = \"http://blog.csdn.net/liang19890820\">快点卖吧</a>")  // 正文

	d := gomail.NewPlainDialer("smtp.qq.com", 25, "947883972@qq.com", "oaeuawfalkcrbfhj")  // 发送邮件服务器、端口、发件人账号、发件人密码
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}