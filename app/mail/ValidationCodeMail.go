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
func NewValidationCodeMail(to, nome, codigo string, isActive bool) *Mail {
	data := ValidationMailData{
		Nome:   nome,
		Codigo: codigo,
		Ano:    time.Now().Year(),
	}

	var tmplPath string

	if !isActive {
		tmplPath = "templates/mails/validation_code_register_mail.html"
	} else {
		tmplPath = "templates/mails/validation_code_change_password_mail.html"
	}

	tmpl, err := template.ParseFiles(tmplPath)
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
		Subject: "Código de Validação - BaseGo",
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
