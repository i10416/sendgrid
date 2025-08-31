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

	// オプションを設定してグローバル統計を取得
	opts := &sendgrid.StatsOptions{
		StartDate:   "2025-01-01",
		EndDate:     "2025-01-31",
		Aggregation: "day",
		Limit:       10,
	}

	stats, err := c.GetGlobalStats(context.TODO(), opts)
	if err != nil {
		return err
	}

	log.Printf("global stats count: %d\n", len(stats))
	for i, stat := range stats {
		log.Printf("stats[%d]: date=%s, delivered=%d, opens=%d, clicks=%d, bounces=%d\n",
			i, stat.Date, stat.Stats.Delivered, stat.Stats.Opens, stat.Stats.Clicks, stat.Stats.Bounces)
	}

	return nil
}
