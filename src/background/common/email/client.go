// Package email is designed to provide an "email interface for humans."
// Designed to be robust and flexible, the email package aims to make sending email easy without getting in the way.
package email

import (
	"net/smtp"
	"time"
)

type MailSender struct {
	smtp      string
	identifer string
	from      string
	username  string
	password  string
	host      string
}

func NewSender(smtp string, host string, identifer string, from string, addr string, pwd string) *MailSender {
	return &MailSender{smtp, identifer, from, addr, pwd, host}
}

func (sender *MailSender) SendMail(to []string, bcc []string, cc []string, subject string, text string, html string) error {
	e := NewEmail()
	e.From = sender.from
	e.To = to
	e.Bcc = bcc
	e.Cc = cc
	e.Subject = subject
	e.Text = []byte(text)
	e.HTML = []byte(html)
	return e.Send(sender.smtp, smtp.PlainAuth(sender.identifer, sender.username, sender.password, sender.host))
}

func (sender *MailSender) SendMailWithAttach(to []string, bcc []string, cc []string, subject string, text string, html string, attachs []string) error {
	e := NewEmail()
	e.From = sender.from
	e.To = to
	e.Bcc = bcc
	e.Cc = cc
	e.Subject = subject
	e.Text = []byte(text)
	e.HTML = []byte(html)
	for _, v := range attachs {
		e.AttachFile(v)
	}
	return e.Send(sender.smtp, smtp.PlainAuth(sender.identifer, sender.username, sender.password, sender.host))
}

func (sender *MailSender) SendBatchMails(mails []*Email) {
	p := NewPool(
		sender.smtp,
		len(mails),
		smtp.PlainAuth(sender.identifer, sender.username, sender.password, sender.host),
	)
	c := make(chan *Email)
	go func() {
		for _, e := range mails {
			c <- e
		}
		close(c)
	}()
	for e := range c {
		p.Send(e, 10*time.Second)
	}
}
