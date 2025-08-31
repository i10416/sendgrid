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
	r, err := c.CreateInboundParseWebhook(context.TODO(), &sendgrid.InputCreateInboundParseWebhook{
		URL:       "https://example.com/sendgrid/inbound",
		Hostname:  "bar.foo",
		SpamCheck: false,
		SendRaw:   false,
	})
	if err != nil {
		return err
	}
	log.Printf("%#v\n", r)

	return nil
}
