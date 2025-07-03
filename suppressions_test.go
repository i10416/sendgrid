package sendgrid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuppressionListOptions(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	opts := &SuppressionListOptions{
		StartTime: 1609459200, // 2021-01-01 00:00:00 UTC
		EndTime:   1612137600, // 2021-02-01 00:00:00 UTC
		Limit:     50,
		Offset:    10,
		Email:     "test@example.com",
	}

	path := "/suppression/bounces"
	fullPath, err := client.AddOptions(path, opts)
	assert.NoError(t, err)
	assert.Contains(t, fullPath, "start_time=1609459200")
	assert.Contains(t, fullPath, "end_time=1612137600")
	assert.Contains(t, fullPath, "limit=50")
	assert.Contains(t, fullPath, "offset=10")
	assert.Contains(t, fullPath, "email=test%40example.com")
}

func TestInputDeleteSuppressions(t *testing.T) {
	// Test with specific emails
	input1 := &InputDeleteSuppressions{
		Emails: []string{"test1@example.com", "test2@example.com"},
	}
	assert.Len(t, input1.Emails, 2)
	assert.False(t, input1.DeleteAll)

	// Test with delete all flag
	input2 := &InputDeleteSuppressions{
		DeleteAll: true,
	}
	assert.True(t, input2.DeleteAll)
	assert.Empty(t, input2.Emails)
}

func TestBounceStruct(t *testing.T) {
	bounce := Bounce{
		Created: 1609459200,
		Email:   "test@example.com",
		Reason:  "550 5.1.1 User unknown",
		Status:  "5.1.1",
	}

	assert.Equal(t, int64(1609459200), bounce.Created)
	assert.Equal(t, "test@example.com", bounce.Email)
	assert.Equal(t, "550 5.1.1 User unknown", bounce.Reason)
	assert.Equal(t, "5.1.1", bounce.Status)
}

func TestBlockStruct(t *testing.T) {
	block := Block{
		Created: 1609459200,
		Email:   "test@example.com",
		Reason:  "IP temporarily blocked",
	}

	assert.Equal(t, int64(1609459200), block.Created)
	assert.Equal(t, "test@example.com", block.Email)
	assert.Equal(t, "IP temporarily blocked", block.Reason)
}

func TestSpamReportStruct(t *testing.T) {
	spamReport := SpamReport{
		Created: 1609459200,
		Email:   "test@example.com",
		IP:      "192.168.1.1",
	}

	assert.Equal(t, int64(1609459200), spamReport.Created)
	assert.Equal(t, "test@example.com", spamReport.Email)
	assert.Equal(t, "192.168.1.1", spamReport.IP)
}

func TestInvalidEmailStruct(t *testing.T) {
	invalidEmail := InvalidEmail{
		Created: 1609459200,
		Email:   "invalid@example.com",
		Reason:  "Mail domain mentioned in email address is unknown",
	}

	assert.Equal(t, int64(1609459200), invalidEmail.Created)
	assert.Equal(t, "invalid@example.com", invalidEmail.Email)
	assert.Equal(t, "Mail domain mentioned in email address is unknown", invalidEmail.Reason)
}
