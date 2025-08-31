package main

import (
	"context"
	"log"
	"os"

	"github.com/i10416/sendgrid"
)

func main() {
	if err := handler(); err != nil {
		log.Fatal(err)
	}
}

func handler() error {
	apiKey := os.Getenv("SENDGRID_API_KEY")

	c := sendgrid.New(apiKey, sendgrid.OptionDebug(true))
	alert, err := c.UpdateAlert(context.TODO(), 8847723, &sendgrid.InputUpdateAlert{
		EmailTo:    "dummy-v2@example.com",
		Frequency:  "weekly",
		Percentage: 99,
	})
	if err != nil {
		return err
	}

	log.Printf("%#v", alert)

	return nil
}
