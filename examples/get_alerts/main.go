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
	alerts, err := c.GetAlerts(context.TODO())
	if err != nil {
		return err
	}

	for _, alert := range alerts {
		log.Printf("alert id: %d type: %s email to: %s\n", alert.ID, alert.Type, alert.EmailTo)
	}

	return nil
}
