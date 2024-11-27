package mail

import (
	"tradeapi/app"
	"os"
	"time"
	"net/smtp"
	"log"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})
}

//enviado o e-mail para a fila
func (m *Mail) SendWithQueue() error {
	logrus.Info("Adicionando e-mail para à fila...")
	
	queue := app.GetQueue()
	queue.ProcessMessages(func(mensagem string) {
		msg := "To: " + m.To + "\r\n" +
			"Subject: " + m.Subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" + mensagem

		err := smtp.SendMail(os.Getenv("MAIL_HOST")+":"+os.Getenv("MAIL_PORT"), 
		smtp.PlainAuth("", os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_PASSWORD"), os.Getenv("MAIL_HOST")), 
		os.Getenv("MAIL_ADDRESS"), []string{m.To}, []byte(msg))
		if err != nil {
			log.Println(err)
		}
	},  1*time.Minute)
	queue.AddMessage(m.Body)
	queue.Wait()
	return nil
}