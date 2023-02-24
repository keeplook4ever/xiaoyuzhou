package mail

import (
	"github.com/jordan-wright/email"
	"net/smtp"
)

func SendMail(to []string, Sub string, txt string) error {
	e := email.NewEmail()
	e.From = "QianYe <support@kiminouchuu.com>"
	e.To = to
	e.Subject = Sub
	e.Text = []byte(txt)
	err := e.Send("smtp.mailgun.org:587", smtp.PlainAuth("", "postmaster@mail.kiminouchuu.com", "45a51a08f1ae8141ee710aeb48cb48c1-52d193a0-8f609649", "smtp.mailgun.org"))
	return err
}
