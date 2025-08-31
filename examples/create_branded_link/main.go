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
	r, err := c.CreateBrandedLink(context.TODO(), &sendgrid.InputCreateBrandedLink{
		Domain:    "examle.com",
		Subdomain: "abc",
	})
	if err != nil {
		return err
	}

	log.Printf("branded link: %#v", r)

	return nil
}
