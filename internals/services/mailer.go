package services

import (
	"time"

	"github.com/wneessen/go-mail"
)

// currently doens't work
type MailerService struct {
	Client *mail.Client
	Config MailerSMTPConfig
}

type MailerSMTPConfig struct {
	Host     string
	Port     int64
	Login    string
	Password string
	Timeout  time.Duration
}

type Mail struct {
	From    string
	To      string
	Subject string
	Content string
}

func NewMailerService(config MailerSMTPConfig) (*MailerService, error) {
	var zeroDuration time.Duration
	if config.Timeout == zeroDuration {
		config.Timeout = 3 * time.Second
	}
	client, err := mail.NewClient(
		config.Host, mail.WithDebugLog(),
		mail.WithPort(int(config.Port)),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(config.Login),
		mail.WithPassword(config.Password),
		mail.WithTimeout(config.Timeout),
	)
	if err != nil {
		return nil, err
	}
	mailer := &MailerService{
		Config: config,
		Client: client,
	}
	return mailer, nil
}

// send mail
func (m MailerService) Perform(input *Mail) error {
	message := mail.NewMsg()
	err := message.From(input.From)
	if err != nil {
		return err
	}
	err = message.To(input.To)
	if err != nil {
		return err
	}
	message.Subject(input.Subject)
	message.SetBodyString(mail.TypeTextPlain, input.Content)

	return m.Client.DialAndSend(message)
}
