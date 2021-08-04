package services

import (
	"crypto/tls"
	"fmt"
	"os"
	"time"

	gomail "gopkg.in/mail.v2"
)

func SendMail(linkAuthor string, authorMail string, originanlLink string, shortUrl string, date time.Time) {
	FROM := os.Getenv("GMAIL")
	PASSWORD := os.Getenv("GMAIL_PASSWORD")
	m := gomail.NewMessage()
	m.SetHeader("From", FROM)
	m.SetHeader("To", authorMail)
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "LinkBook")
	mailBody := "Hello , " + "<b>" + linkAuthor + "</b>" + "&#9825;" + ".<br><br>" + "Your Link - " + originanlLink + " has been Published . " + "<br>" + " Short Url :- " + shortUrl + "<br>" + " Time - " + date.String() + "<br><br>" + "Thanks!"
	m.SetBody("text/html", mailBody)
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.gmail.com", 587, FROM, PASSWORD)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	} else {
		fmt.Println("Mail sent Done")
	}
}
