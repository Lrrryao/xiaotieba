package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

// 定义一个邮件发送者的接口
// 发送者需要完成SendEmail函数才能注册接口
type EmailSender interface {
	SendEmail(subject string, content string, to, cc, bcc, attachFiles []string) error
} //标题subject， 内容content，接收人to，抄送cc，暗抄送bcc，附件attachFiles

type QQSender struct {
	senderName        string //发件人name
	fromEmailAddress  string //发件人地址
	fromEmailPassword string //发件人邮箱密码，不会使用真实的密码，
}

func NewQQSender(name, addr, password string) EmailSender {

	return &QQSender{
		senderName:        name,
		fromEmailAddress:  addr,
		fromEmailPassword: password,
	}
}

const (
	smtpAuthAddress   = "smtp.qq.com"
	smtpServerAddress = "smtp.qq.com:465"
)

func (sender *QQSender) SendEmail(subject string, content string, to, cc, bcc, attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.senderName, sender.fromEmailAddress)
	e.To = to
	e.Bcc = bcc
	e.Cc = cc
	e.Subject = subject
	e.Text = []byte("welcome to xaiotieba")
	e.HTML = []byte(content)

	for _, file := range attachFiles {
		_, err := e.AttachFile(file)

		if err != nil {
			return fmt.Errorf("fail to attach file %s, err:%v", file, err)
		}
	}
	auth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, auth)

}
