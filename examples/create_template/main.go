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
	r, err := c.CreateTemplate(context.TODO(), &sendgrid.InputCreateTemplate{
		Name:       "dummy",
		Generation: "dynamic",
	})
	if err != nil {
		return err
	}

	log.Printf("temlate: %#v", r)

	return nil
}
