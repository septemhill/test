package utils

import "gopkg.in/gomail.v2"

type MailType int

const (
	MAIL_TYPE_HTML = iota
	MAIL_TYPE_PLAIN
)

var mailTypeMap = map[MailType]string{
	MAIL_TYPE_HTML:  "text/html",
	MAIL_TYPE_PLAIN: "text/plain",
}

type MailInfo struct {
	From    string
	To      string
	Subject string
	Body    string
	Type    MailType
}

type Mailer struct {
	Host     string `env:"MAILER_SERVER_HOST"`
	Port     int    `env:"MAILER_SERVER_SMTP_PORT" envDefault:"587"`
	User     string `env:"MAILER_ID"`
	Password string `env:"MAILER_PASSWORD"`
}

func SendMail(mailer Mailer, info MailInfo) error {
	m := gomail.NewMessage()
	m.SetHeader("From", info.From)
	m.SetHeader("To", info.To)
	m.SetHeader("Subject", info.Subject)
	m.SetBody(mailTypeMap[info.Type], info.Body)

	return gomail.NewDialer(
		mailer.Host,
		mailer.Port,
		mailer.User,
		mailer.Password,
	).DialAndSend(m)
}
