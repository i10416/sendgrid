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

	// オプションを設定してバウンス一覧を取得
	opts := &sendgrid.SuppressionListOptions{
		Limit:  10,
		Offset: 0,
	}

	bounces, err := c.GetBounces(context.TODO(), opts)
	if err != nil {
		return err
	}

	log.Printf("bounces count: %d\n", len(bounces))
	for i, bounce := range bounces {
		log.Printf("bounce[%d]: email=%s, reason=%s, status=%s\n",
			i, bounce.Email, bounce.Reason, bounce.Status)
	}

	return nil
}
