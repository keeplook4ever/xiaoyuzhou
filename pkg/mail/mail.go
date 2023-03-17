package mail

import (
	"github.com/jordan-wright/email"
	"net/smtp"
)

func SendMail(to []string, Sub string, htmlTxt string) error {
	e := email.NewEmail()
	e.From = "business@kiminouchuu.com"
	e.To = to
	e.Subject = Sub
	e.Text = []byte("ssssss")
	err := e.Send("smtp.mailgun.org:587", smtp.PlainAuth("", "business@kiminouchuu.com", "f35e2f4d952f5d3b8b951cb596b5ce92-b36d2969-ef33b4cc", "smtp.mailgun.org"))
	return err
}
