package helpers

import (
	"fmt"
	"net/http"
	"time"

	"gopkg.in/gomail.v2"
)

func tenantsendEmail(subject, message, toAddress string) error {
	// Create a new message
	m := gomail.NewMessage()
	m.SetHeader("From", "your-email@example.com") // Replace with your email address
	m.SetHeader("To", toAddress)                  // Recipient's email address
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", message)

	// Create a new dialer to establish a connection to the SMTP server
	d := gomail.NewDialer("smtp.example.com", 587, "your-email@example.com", "your-password") // Replace with your SMTP server details

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func CheckURLAndSendEmail(tenantURL string) {
	url := tenantURL // Replace with your tenant's URL

	for {
		// Check the status of the tenant's URL
		response, err := http.Get(url)
		if err != nil {
			// Handle the error, you might want to log it
			fmt.Println("Failed to check tenant URL status:", err)
		} else if response.StatusCode == http.StatusOK {
			// If the URL responds with a status code of 200, send an email
			emailErr := tenantsendEmail("Deployment Successful", "Your tenant has been deployed successfully.", "recipient@example.com")
			if emailErr != nil {
				// Handle the email sending error, you might want to log it
				fmt.Println("Failed to send email:", emailErr)
			}
		}
		time.Sleep(15 * time.Second) // Check every 15 seconds
	}
}
