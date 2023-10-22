package services

import "log"

type Mailer struct {
	From    string
	To      string
	Subject string
	Content string
}

func (m Mailer) SendMailer() error {
	log.Println("Mail sent")
	return nil
}
