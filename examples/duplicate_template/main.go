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
	r, err := c.DuplicateTemplate(context.TODO(), "d-12345abcde", &sendgrid.InputDuplicateTemplate{
		Name: "dummy2",
	})
	if err != nil {
		return err
	}

	log.Printf("temlate: %#v", r)

	return nil
}
