package sendgrid

import (
	"context"
	"fmt"
)

type OutputGetAlert struct {
	ID         int64  `json:"id,omitempty"`
	EmailTo    string `json:"email_to,omitempty"`
	Frequency  string `json:"frequency,omitempty"`
	Type       string `json:"type,omitempty"`
	Percentage int64  `json:"percentage,omitempty"`
	CreatedAt  int64  `json:"created_at,omitempty"`
	UpdatedAt  int64  `json:"updated_at,omitempty"`
}

// see: https://www.twilio.com/docs/sendgrid/api-reference/alerts/retrieve-a-specific-alert
func (c *Client) GetAlert(ctx context.Context, id int64) (*OutputGetAlert, error) {
	path := fmt.Sprintf("/alerts/%d", id)

	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	r := new(OutputGetAlert)
	if err := c.Do(ctx, req, &r); err != nil {
		return nil, err
	}

	return r, nil
}

type Alert struct {
	ID         int64  `json:"id,omitempty"`
	EmailTo    string `json:"email_to,omitempty"`
	Frequency  string `json:"frequency,omitempty"`
	Type       string `json:"type,omitempty"`
	Percentage int64  `json:"percentage,omitempty"`
	CreatedAt  int64  `json:"created_at,omitempty"`
	UpdatedAt  int64  `json:"updated_at,omitempty"`
}

// see: https://www.twilio.com/docs/sendgrid/api-reference/alerts/retrieve-all-alerts
func (c *Client) GetAlerts(ctx context.Context) ([]*Alert, error) {
	req, err := c.NewRequest("GET", "/alerts", nil)
	if err != nil {
		return nil, err
	}

	r := []*Alert{}
	if err := c.Do(ctx, req, &r); err != nil {
		return nil, err
	}

	return r, nil
}

type InputCreateAlert struct {
	Type       string `json:"type,omitempty"`
	EmailTo    string `json:"email_to,omitempty"`
	Frequency  string `json:"frequency,omitempty"`
	Percentage int64  `json:"percentage,omitempty"`
}

type OutputCreateAlert struct {
	ID         int64  `json:"id,omitempty"`
	EmailTo    string `json:"email_to,omitempty"`
	Frequency  string `json:"frequency,omitempty"`
	Type       string `json:"type,omitempty"`
	Percentage int64  `json:"percentage,omitempty"`
	CreatedAt  int64  `json:"created_at,omitempty"`
	UpdatedAt  int64  `json:"updated_at,omitempty"`
}

// see: https://www.twilio.com/docs/sendgrid/api-reference/alerts/create-a-new-alert
func (c *Client) CreateAlert(ctx context.Context, input *InputCreateAlert) (*OutputCreateAlert, error) {
	req, err := c.NewRequest("POST", "/alerts", input)
	if err != nil {
		return nil, err
	}

	r := new(OutputCreateAlert)
	if err := c.Do(ctx, req, &r); err != nil {
		return nil, err
	}

	return r, nil
}

type InputUpdateAlert struct {
	EmailTo    string `json:"email_to,omitempty"`
	Frequency  string `json:"frequency,omitempty"`
	Percentage int64  `json:"percentage,omitempty"`
}

type OutputUpdateAlert struct {
	ID         int64  `json:"id,omitempty"`
	EmailTo    string `json:"email_to,omitempty"`
	Frequency  string `json:"frequency,omitempty"`
	Type       string `json:"type,omitempty"`
	Percentage int64  `json:"percentage,omitempty"`
	CreatedAt  int64  `json:"created_at,omitempty"`
	UpdatedAt  int64  `json:"updated_at,omitempty"`
}

// see: https://www.twilio.com/docs/sendgrid/api-reference/alerts/update-an-alert
func (c *Client) UpdateAlert(ctx context.Context, id int64, input *InputUpdateAlert) (*OutputUpdateAlert, error) {
	path := fmt.Sprintf("/alerts/%d", id)

	req, err := c.NewRequest("PATCH", path, input)
	if err != nil {
		return nil, err
	}

	r := new(OutputUpdateAlert)
	if err := c.Do(ctx, req, &r); err != nil {
		return nil, err
	}

	return r, nil
}

// see: https://www.twilio.com/docs/sendgrid/api-reference/alerts/delete-an-alert
func (c *Client) DeleteAlert(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/alerts/%d", id)

	req, err := c.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	if err := c.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}
