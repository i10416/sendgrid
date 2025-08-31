package sendgrid

import (
	"context"
	"fmt"
	"net/url"
)

type CustomField struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type InputGetCustomField struct {
	ID int64
}

func (c *Client) GetCustomField(ctx context.Context, id int64) (*CustomField, error) {
	u, _ := url.Parse(fmt.Sprintf("/contactdb/custom_fields/%d", id))

	req, err := c.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	r := CustomField{}
	if err := c.Do(ctx, req, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

type InputCreateCustomField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (c *Client) CreateCustomField(ctx context.Context, input *InputCreateCustomField) (*CustomField, error) {
	req, err := c.NewRequest("POST", "/contactdb/custom_fields", input)
	if err != nil {
		return nil, err
	}

	r := new(CustomField)
	if err := c.Do(ctx, req, &r); err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Client) DeleteCustomField(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/contactdb/custom_fields/%d", id)

	req, err := c.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	if err := c.Do(ctx, req, nil); err != nil {
		return err
	}
	return nil
}
