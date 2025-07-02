package sendgrid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatsOptions(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	opts := &StatsOptions{
		StartDate:   "2025-01-01",
		EndDate:     "2025-01-31",
		Aggregation: "day",
		Limit:       100,
		Offset:      10,
	}

	path := "/stats"
	fullPath, err := client.AddOptions(path, opts)
	assert.NoError(t, err)
	assert.Contains(t, fullPath, "start_date=2025-01-01")
	assert.Contains(t, fullPath, "end_date=2025-01-31")
	assert.Contains(t, fullPath, "aggregated_by=day")
	assert.Contains(t, fullPath, "limit=100")
	assert.Contains(t, fullPath, "offset=10")
}

func TestStatMetrics(t *testing.T) {
	metrics := StatMetrics{
		Blocks:           5,
		BounceDrops:      2,
		Bounces:          3,
		Clicks:           100,
		DeferredDrops:    1,
		Delivered:        950,
		InvalidEmails:    4,
		Opens:            300,
		Processed:        1000,
		Requests:         1000,
		SpamReportDrops:  1,
		SpamReports:      2,
		UniqueClicks:     85,
		UniqueOpens:      250,
		UnsubscribeDrops: 1,
		Unsubscribes:     3,
	}

	assert.Equal(t, 5, metrics.Blocks)
	assert.Equal(t, 2, metrics.BounceDrops)
	assert.Equal(t, 3, metrics.Bounces)
	assert.Equal(t, 100, metrics.Clicks)
	assert.Equal(t, 1, metrics.DeferredDrops)
	assert.Equal(t, 950, metrics.Delivered)
	assert.Equal(t, 4, metrics.InvalidEmails)
	assert.Equal(t, 300, metrics.Opens)
	assert.Equal(t, 1000, metrics.Processed)
	assert.Equal(t, 1000, metrics.Requests)
	assert.Equal(t, 1, metrics.SpamReportDrops)
	assert.Equal(t, 2, metrics.SpamReports)
	assert.Equal(t, 85, metrics.UniqueClicks)
	assert.Equal(t, 250, metrics.UniqueOpens)
	assert.Equal(t, 1, metrics.UnsubscribeDrops)
	assert.Equal(t, 3, metrics.Unsubscribes)
}

func TestGlobalStat(t *testing.T) {
	stat := GlobalStat{
		Date: "2025-01-01",
		Stats: StatMetrics{
			Delivered: 1000,
			Opens:     300,
			Clicks:    100,
		},
	}

	assert.Equal(t, "2025-01-01", stat.Date)
	assert.Equal(t, 1000, stat.Stats.Delivered)
	assert.Equal(t, 300, stat.Stats.Opens)
	assert.Equal(t, 100, stat.Stats.Clicks)
}

func TestStatItem(t *testing.T) {
	item := StatItem{
		Name: "newsletter",
		Type: "category",
		Metrics: StatMetrics{
			Delivered: 500,
			Opens:     150,
			Clicks:    50,
		},
	}

	assert.Equal(t, "newsletter", item.Name)
	assert.Equal(t, "category", item.Type)
	assert.Equal(t, 500, item.Metrics.Delivered)
	assert.Equal(t, 150, item.Metrics.Opens)
	assert.Equal(t, 50, item.Metrics.Clicks)
}