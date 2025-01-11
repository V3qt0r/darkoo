package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	mailgun "github.com/mailgun/mailgun-go/v4"
)


func SendEmail(sender, recipient, subject, body string) error {

	// Create an instance of the Mailgun Client
	mg, _ := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(mailgun.APIBaseEU)

	//When you have an EU-domain, you must specify the endpoint:
	//mg.SetAPIBase("https://api.eu.mailgun.net/v3")

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, "", recipient)
	message.SetHtml(body)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Print(err)
		return err
	}
	
	fmt.Printf("ID: %s resp: %s\n", id, resp)
	return nil
}


func SendEmailWithDefaultSender(recipient, subject, body string) error{
	sender := os.Getenv("EMAIL_SENDER_EMAIL")

	if len(strings.TrimSpace(sender)) == 0 {
		return errors.New("Sender email not configured")
	}

	err := SendEmail(sender, recipient, subject, body)
	if err != nil {
		return err
	}

	return nil
}