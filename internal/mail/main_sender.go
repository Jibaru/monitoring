package mail

import (
	"net/smtp"
)

type MailSender struct {
	from            string
	mailAppPassword string
}

func NewMailSender(from, mailAppPassword string) *MailSender {
	return &MailSender{
		from:            from,
		mailAppPassword: mailAppPassword,
	}
}

func (m *MailSender) Send(to, subject, body string) error {
	msg := "From: " + m.from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", m.from, m.mailAppPassword, "smtp.gmail.com"),
		m.from, []string{to}, []byte(msg))

	if err != nil {
		return err
	}

	return nil
}
