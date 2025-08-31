package sendgrid

import (
	"context"
	"fmt"
	"net/url"
)

type Segment struct {
	ID         int64              `json:"id,omitempty"`
	Name       string             `json:"name"`
	Conditions []SegmentCondition `json:"conditions"`
}

func (c *Client) GetSegment(ctx context.Context, id int64) (*Segment, error) {
	u, _ := url.Parse(fmt.Sprintf("/contactdb/segments/%d", id))

	req, err := c.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	r := Segment{}
	if err := c.Do(ctx, req, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

type InputCreateSegment struct {
	Name       string             `json:"name"`
	Conditions []SegmentCondition `json:"conditions"`
}

type InputUpdateSegment struct {
	Name       string             `json:"name,omitempty"`
	Conditions []SegmentCondition `json:"conditions,omitempty"`
}

type SegmentCondition struct {
	Field    string `json:"field"`
	Value    string `json:"value"`
	Operator string `json:"operator"`
	AndOr    string `json:"and_or"`
}

func (c *Client) CreateSegment(ctx context.Context, input *InputCreateSegment) (*Segment, error) {
	req, err := c.NewRequest("POST", "/contactdb/segments", input)
	if err != nil {
		return nil, err
	}

	r := new(Segment)
	if err := c.Do(ctx, req, &r); err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Client) UpdateSegment(ctx context.Context, input *InputUpdateSegment) (*Segment, error) {
	req, err := c.NewRequest("POST", "/contactdb/segments", input)
	if err != nil {
		return nil, err
	}

	r := new(Segment)
	if err := c.Do(ctx, req, &r); err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Client) DeleteSegment(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/contactdb/segments/%d", id)

	req, err := c.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	if err := c.Do(ctx, req, nil); err != nil {
		return err
	}
	return nil
}
