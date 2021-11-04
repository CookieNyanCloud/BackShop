package smtp

import (
	"fmt"
	"github.com/cookienyancloud/back/pkg/email"
	"github.com/pkg/errors"
	"net/smtp"
	"strconv"
)

type Sender struct {
	from string
	pass string
	host string
	port int
}

func NewSMTPSender(from string, pass string, host string, port int) (*Sender, error) {
	if !email.IsEmailValid(from) {
		return nil, errors.New("invalid from email")
	}

	return &Sender{from: from, pass: pass, host: host, port: port}, nil
}

func (s *Sender) Send(input email.SendEmailInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	//msg := gomail.NewMessage()
	//msg.SetHeader("From", s.from)
	//msg.SetHeader("To", input.To)
	//msg.SetHeader("Subject", input.Subject)
	//msg.SetBody("text/html", input.Body)

	auth := smtp.PlainAuth("", s.from, s.pass, s.host)
	address := s.host + ":" + strconv.Itoa(s.port)
	message := []byte(input.Subject + input.Body)
	to := []string{input.To}
	fmt.Println("message:", string(message))
	err := smtp.SendMail(address, auth, s.from, to, message)
	if err != nil {
		return errors.Wrap(err,"cannot send message")
	}
	return nil
}
