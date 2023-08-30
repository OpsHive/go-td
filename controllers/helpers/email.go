package helpers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

var (
	Domain string = os.Getenv("MAIN_GUN_DOMAIN") //eg tenanatmail.tripon.io

	PrivateAPIKey string = os.Getenv("MAIL_GUN_API_KEY")
)

func CheckTenantUrlSendEmail(tenantURL, toAddress, subject, message string, sleepInterval time.Duration) error {
	for {
		// Check the status of the tenant's URL
		response, err := http.Get(tenantURL)
		if err != nil {
			// Handle the error, you might want to return or log it
			return err
		} else if response.StatusCode == http.StatusOK {
			// If the URL responds with a status code of 200, send an email
			emailErr := sendEmail(subject, message, toAddress)
			if emailErr != nil {
				// Handle the email sending error, you might want to return or log it
				return emailErr
			}
		}
		time.Sleep(sleepInterval)
	}
}

func sendEmail(subject, message, toAddress string) error {

	mg := mailgun.NewMailgun(Domain, PrivateAPIKey)

	sender := os.Getenv("MAIL_GUN_SENDER")
	subj := subject
	body := message
	recipient := toAddress

	m := mg.NewMessage(sender, subj, body, recipient)

	// Send the email without a context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, m)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
	return err
}
