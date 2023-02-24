package mail

import (
	"context"
	"github.com/mailgun/mailgun-go/v3"
	"time"
)

func SendSimpleMessage(domain, apiKey string) (string, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage(
		"Admin <cw4t@mail.kiminouchuu.com>",
		"Hello",
		"Testing some Mailgun awesomeness!",
		"support@mail.kiminouchuu.com", "cwater4t@gmail.com",
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}
