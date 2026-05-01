package mailer

import (
	"crypto/tls"
	"fmt"

	"github.com/WJH9102/msg_forwarder/internal/config"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Mailer {
	return &Mailer{cfg: cfg}
}

func (m *Mailer) Send(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", fmt.Sprintf("%s <%s>", m.cfg.SenderName, m.cfg.SMTPUser))
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	dialer := gomail.NewDialer(
		m.cfg.SMTPHost,
		m.cfg.SMTPPort,
		m.cfg.SMTPUser,
		m.cfg.SMTPPassword,
	)
	dialer.TLSConfig = &tls.Config{ServerName: m.cfg.SMTPHost}

	return dialer.DialAndSend(msg)
}