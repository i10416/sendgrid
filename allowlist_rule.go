package sendgrid

import (
	"context"
	"fmt"
	"net/url"
)

type OutputGetAllowlistRule struct {
	Result AllowlistRule `json:"result"`
}

type AllowlistRule struct {
	ID int64  `json:"id"`
	Ip string `json:"ip"`
}

func (c *Client) GetAllowlistRule(ctx context.Context, id int64) (*AllowlistRule, error) {
	u, _ := url.Parse(fmt.Sprintf("/access_settings/whitelist/%d", id))

	req, err := c.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	r := OutputGetAllowlistRule{}
	if err := c.Do(ctx, req, &r); err != nil {
		return nil, err
	}

	return &r.Result, nil
}

type InputCreateAllowlistRuleIp struct {
	Ip string `json:"ip"`
}

type InputCreateAllowlistRule struct {
	Ips []InputCreateAllowlistRuleIp `json:"ips"`
}
type OutputCreateAllowlistRule struct {
	Result []AllowlistRule `json:"result"`
}

func (c *Client) CreateAllowlistRule(ctx context.Context, input *InputCreateAllowlistRule) (*OutputCreateAllowlistRule, error) {
	req, err := c.NewRequest("POST", "/access_settings/whitelist", input)
	if err != nil {
		return nil, err
	}

	r := new(OutputCreateAllowlistRule)
	if err := c.Do(ctx, req, &r); err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Client) DeleteAllowlistRule(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/access_settings/whitelist/%d", id)

	req, err := c.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	if err := c.Do(ctx, req, nil); err != nil {
		return err
	}
	return nil
}
