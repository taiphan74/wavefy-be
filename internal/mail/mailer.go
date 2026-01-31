package mail

import (
	"crypto/tls"
	"errors"
	"fmt"

	"gopkg.in/gomail.v2"

	"wavefy-be/config"
)

type Config struct {
	Host string
	Port int
	User string
	Pass string
	From string
}

type Service struct {
	dialer *gomail.Dialer
	from   string
}

func New(cfg Config) (*Service, error) {
	if cfg.Host == "" || cfg.Port == 0 || cfg.User == "" || cfg.Pass == "" || cfg.From == "" {
		return nil, errors.New("smtp config is incomplete")
	}

	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Pass)
	d.TLSConfig = &tls.Config{
		ServerName: cfg.Host,
		MinVersion: tls.VersionTLS12,
	}

	return &Service{dialer: d, from: cfg.From}, nil
}

func FromConfig(cfg config.MailConfig) (*Service, error) {
	if cfg.Host == "" || cfg.Port == 0 || cfg.User == "" || cfg.Pass == "" || cfg.From == "" {
		return nil, nil
	}
	return New(Config{
		Host: cfg.Host,
		Port: cfg.Port,
		User: cfg.User,
		Pass: cfg.Pass,
		From: cfg.From,
	})
}

func (s *Service) Send(to, subject, textBody, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	if htmlBody != "" {
		m.SetBody("text/html", htmlBody)
		if textBody != "" {
			m.AddAlternative("text/plain", textBody)
		}
	} else {
		m.SetBody("text/plain", textBody)
	}

	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("send mail failed: %w", err)
	}
	return nil
}
