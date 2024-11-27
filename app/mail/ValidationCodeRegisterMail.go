package mail

import (
	"bytes"
	"html/template"
	"time"
)


type Mail struct {
	To      string
	Subject string
	Body    string
}

type ValidationMailData struct {
	Nome   string
	Codigo string
	Ano    int
}

// NewValidationCodeMail cria um novo e-mail de código de validação
func NewValidationCodeMail(to, nome, codigo string) *Mail {
	data := ValidationMailData{
		Nome:   nome,
		Codigo: codigo,
		Ano:    time.Now().Year(),
	}

	tmpl, err := template.ParseFiles("templates/mails/validation_code_register_mail.html")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}

	mail := &Mail{
		To:      to,
		Subject: "Código de Validação - Tradeapi",
		Body:    buf.String(),
	}

	mail.SendWithQueue()

	return mail
}

// Send envia o e-mail diretamente
/*func (m *Mail) Send() error {
	auth := smtp.PlainAuth("", os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_PASSWORD"), os.Getenv("MAIL_EMAIL"))

	msg := "To: " + m.To + "\r\n" +
		"Subject: " + m.Subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + m.Body

	return smtp.SendMail("smtp.mailtrap.io:2525", auth, os.Getenv("MAIL_EMAIL"), []string{m.To}, []byte(msg))
}*/
