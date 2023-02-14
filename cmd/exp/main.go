package main

import (
	"log"
	"os"

	"github.com/go-mail/mail/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	from := "test@lenslocked.com"
	to := "snir@example.com"
	subject := "This is a test email"
	plaintext := "Body of the email"
	html := `<h1>Hello there budyy!</h1><p>This is the amil</p>`

	msg := mail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetHeader("From", from)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plaintext)
	msg.AddAlternative("text/html", html)

	msg.WriteTo(os.Stdout)

	d := mail.NewDialer("sandbox.smtp.mailtrap.io", 25, "", "")

	if err := d.DialAndSend(msg); err != nil {
		log.Fatal(err)
	}
}
