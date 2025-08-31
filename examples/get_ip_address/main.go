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
	ip, err := c.GetIPAddress(context.TODO(), "192.168.1.1")
	if err != nil {
		return err
	}
	log.Printf("ip address: %#v\n", ip)
	return nil
}
