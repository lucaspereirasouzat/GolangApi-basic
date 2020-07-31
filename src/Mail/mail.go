// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mail

import (
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func Example(mail, message string) {
	// Connect to the remote SMTP server.
	fmt.Println(os.Getenv("SMTP_SERVER") + ":" + os.Getenv("SMTP_PORT"))

	conn, err := net.Dial("tcp", os.Getenv("SMTP_SERVER")+":"+os.Getenv("SMTP_PORT"))
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println(conn)
	smtpAdress := os.Getenv("SMTP_SERVER") + ":" + os.Getenv("SMTP_PORT")
	c, err := smtp.Dial(smtpAdress)
	fmt.Println(c, err)
	if err != nil {
		log.Fatal(err)
	}

	// Set the sender and recipient first
	if err := c.Mail(os.Getenv("EMAIL_SMTP")); err != nil {
		log.Fatal(err)
	}
	if err := c.Rcpt(mail); err != nil {
		log.Fatal(err)
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintf(wc, message)
	if err != nil {
		log.Fatal(err)
	}
	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		log.Fatal(err)
	}
}

// // variables to make ExamplePlainAuth compile, without adding
// // unnecessary noise there.
// var (
// 	from       = "gopher@example.net"
// 	msg        = []byte("dummy message")
// 	recipients = []string{"foo@example.com"}
// )

// func ExamplePlainAuth() {
// 	// hostname is used by PlainAuth to validate the TLS certificate.
// 	hostname := "mail.example.com"
// 	auth := smtp.PlainAuth("", "user@example.com", "password", hostname)

// 	err := smtp.SendMail(hostname+":25", auth, from, recipients, msg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func ExampleSendMail(to string, msg string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o .env")
	}

	msgs := []byte(msg)
	// Set up authentication information.
	auth := smtp.PlainAuth("", os.Getenv("EMAIL_SMTP"), os.Getenv("EMAIL_Password"), os.Getenv("SMTP_SERVER"))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	toB := []string{to}
	// msgB := []byte(to + "\r\n" +
	// 	"Subject: discount Gophers!\r\n" +
	// 	"\r\n" +
	// 	"This is the email body.\r\n")
	fmt.Println(msgs, auth, toB)
	err = smtp.SendMail(os.Getenv("SMTP_SERVER")+":"+os.Getenv("SMTP_PORT"), auth, os.Getenv("EMAIL_SMTP"), toB, msgs)
	if err != nil {
		log.Fatal(err)
	}
}
