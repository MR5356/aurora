package notify

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/MR5356/aurora/pkg/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type EmailNotifier struct {
	dialer *gomail.Dialer
}

func NewEmailNotifier(conf config.Email) *EmailNotifier {
	logrus.Infof("new email notifier: %+v", conf)
	dialer := gomail.NewDialer(conf.Host, conf.Port, conf.Username, conf.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &EmailNotifier{
		dialer: dialer,
	}
}

func (n *EmailNotifier) Send(ctx context.Context, msg *MessageTemplate) error {
	logrus.Infof("send email: %+v", msg)
	m := gomail.NewMessage()
	m.SetHeader("From", n.dialer.Username)
	m.SetHeader("From", config.Current().Email.Alias+"<"+n.dialer.Username+">")

	m.SetHeader("To", msg.Receivers.Receivers...)
	m.SetHeader("Subject", msg.Subject)

	m.SetBody("text/html", fmt.Sprintf(msg.Body))

	return n.dialer.DialAndSend(m)
}
