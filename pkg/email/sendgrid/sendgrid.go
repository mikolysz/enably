package sendgrid

import (
	"fmt"
	"log"

	"github.com/mikolysz/enably/pkg/email"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Config struct {
	SenderEmail string
	SenderName  string
	APIKey      string
}

type Sender struct {
	config Config
}

func NewSender(config Config) *Sender {
	return &Sender{config: config}
}

func (s *Sender) Send(m email.Message) error {
	from := mail.NewEmail(s.config.SenderName, s.config.SenderEmail)
	to := mail.NewEmail("", m.Recipient)
	msg := mail.NewSingleEmail(from, m.Subject, to, m.PlainTextContent, m.HTMLContent)
	client := sendgrid.NewSendClient(s.config.APIKey)
	resp, err := client.Send(msg)
	if err != nil {
		return fmt.Errorf("error when sending email: %w", err)
	}
	log.Printf("sendgrid response: %+v", resp)
	return nil
}
