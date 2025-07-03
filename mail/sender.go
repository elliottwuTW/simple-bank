package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress = "smtp.gmail.com"
	// Gmail 用來傳送的 address
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	mailAccountName   string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(
	name string,
	fromEmailAddress string,
	fromEmailPassword string,
) EmailSender {
	return &GmailSender{
		mailAccountName:   name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.mailAccountName, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, a := range attachFiles {
		_, err := e.AttachFile(a)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", a, err)
		}
	}

	// 向 SMTP server 進行認證
	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	// 傳送 email
	return e.Send(smtpServerAddress, smtpAuth)
}
