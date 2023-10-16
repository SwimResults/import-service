package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"testing"
	"text/template"
)

func TestMails(tst *testing.T) {
	// Sender data.
	from := "no-reply@swimresults.de"

	user := "swimresults.de"
	password := "dmnzhszqabkmibqc"

	// Receiver email address.
	to := []string{
		"konrad@schwimmteamerzgebirge.de",
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", user, password, smtpHost)

	t, _ := template.ParseFiles("template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    "Konrad Weiß",
		Message: "This is a test message in a HTML template from SwimResults",
	})

	if os.Getenv("SR_SEND_MAILS") == "" {
		return
	}

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}
