package main

import (
	"context"
	"log"
	"os"

	"github.com/kenzo0107/sendgrid"
)

func main() {
	if err := handler(); err != nil {
		log.Fatal(err)
	}
}

func handler() error {
	apiKey := os.Getenv("SENDGRID_API_KEY")

	c := sendgrid.New(apiKey, sendgrid.OptionDebug(true))
	alert, err := c.CreateAlert(context.TODO(), &sendgrid.InputCreateAlert{
		Type:       "stats_notification",
		EmailTo:    "dummy@example.com",
		Frequency:  "daily",
		Percentage: 90,
	})
	if err != nil {
		return err
	}

	log.Printf("%#v", alert)

	return nil
}
