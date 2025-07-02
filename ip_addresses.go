package sendgrid

import (
	"context"
	"fmt"
	"net/url"
)

// IPAddress represents an IP address
type IPAddress struct {
	IP         string   `json:"ip,omitempty"`
	Pools      []string `json:"pools,omitempty"`
	Warmup     bool     `json:"warmup,omitempty"`
	StartDate  int64    `json:"start_date,omitempty"`
	Subusers   []string `json:"subusers,omitempty"`
	Rdns       string   `json:"rdns,omitempty"`
	AssignedAt int64    `json:"assigned_at,omitempty"`
}

// IPPool represents an IP pool
type IPPool struct {
	Name string `json:"name,omitempty"`
}

// IPWarmupStatus represents IP warmup status
type IPWarmupStatus struct {
	IP     string `json:"ip,omitempty"`
	Warmup bool   `json:"warmup,omitempty"`
}

// InputAddIPToPool represents the request to add an IP to a pool
type InputAddIPToPool struct {
	IP string `json:"ip"`
}

// InputAssignIPToSubuser represents the request to assign an IP to a subuser
type InputAssignIPToSubuser struct {
	IPs []string `json:"ips"`
}

// OutputAssignedIPs represents the response for assigned IPs
type OutputAssignedIPs struct {
	IPs []string `json:"ips,omitempty"`
}

// GetIPAddresses retrieves all IP addresses
// see: https://www.twilio.com/docs/sendgrid/api-reference/ip-addresses/retrieve-all-ip-addresses
func (c *Client) GetIPAddresses(ctx context.Context) ([]IPAddress, error) {
	path := "/ips"

	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var ips []IPAddress
	if err := c.Do(ctx, req, &ips); err != nil {
		return nil, err
	}

	return ips, nil
}

// GetAssignedIPAddresses retrieves all assigned IP addresses
// see: https://www.twilio.com/docs/sendgrid/api-reference/ip-addresses/retrieve-all-assigned-ips
func (c *Client) GetAssignedIPAddresses(ctx context.Context) ([]IPAddress, error) {
	path := "/ips/assigned"

	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var ips []IPAddress
	if err := c.Do(ctx, req, &ips); err != nil {
		return nil, err
	}

	return ips, nil
}

// GetRemainingIPCount retrieves remaining IP count
// see: https://www.twilio.com/docs/sendgrid/api-reference/ip-addresses/retrieve-remaining-ip-count
func (c *Client) GetRemainingIPCount(ctx context.Context) (map[string]interface{}, error) {
	path := "/ips/remaining"

	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := c.Do(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetIPAddress retrieves a specific IP address
// see: https://www.twilio.com/docs/sendgrid/api-reference/ip-addresses/retrieve-an-ip-address
func (c *Client) GetIPAddress(ctx context.Context, ip string) (*IPAddress, error) {
	path := fmt.Sprintf("/ips/%s", url.QueryEscape(ip))

	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var ipAddr IPAddress
	if err := c.Do(ctx, req, &ipAddr); err != nil {
		return nil, err
	}

	return &ipAddr, nil
}

// GetIPPools retrieves all IP pools
// see: https://www.twilio.com/docs/sendgrid/api-reference/ip-pools/retrieve-all-ip-pools
func (c *Client) GetIPPools(ctx context.Context) ([]IPPool, error) {
	path := "/ips/pools"

	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var pools []IPPool
	if err := c.Do(ctx, req, &pools); err != nil {
		return nil, err
	}

	return pools, nil
}

// CreateIPPool creates an IP pool
// see: https://www.twilio.com/docs/sendgrid/api-reference/ip-pools/create-an-ip-pool
func (c *Client) CreateIPPool(ctx context.Context, name string) (*IPPool, error) {
	path := "/ips/pools"

	input := &IPPool{Name: name}

	req, err := c.NewRequest("POST", path, input)
	if err != nil {
		return nil, err
	}

	var pool IPPool
	if err := c.Do(ctx, req, &pool); err != nil {
		return nil, err
	}

	return &pool, nil
}

// GetIPPool retrieves a specific IP pool
// see: https://www.twilio.com/docs/sendgrid/api-reference/ip-pools/retrieve-an-ip-pool
func (c *Client) GetIPPool(ctx context.Context, name string) (*IPPool, error) {
	path := fmt.Sprintf("/ips/pools/%s", url.QueryEscape(name))

	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var pool IPPool
	if err := c.Do(ctx, req, &pool); err != nil {
		return nil, err
	}

	return &pool, nil
}

// UpdateIPPool updates an IP pool name
// see: https://www.twilio.com/docs/sendgrid/api-reference/ip-pools/update-an-ip-pool
func (c *Client) UpdateIPPool(ctx context.Context, oldName, newName string) (*IPPool, error) {
	path := fmt.Sprintf("/ips/pools/%s", url.QueryEscape(oldName))

	input := &IPPool{Name: newName}

	req, err := c.NewRequest("PUT", path, input)
	if err != nil {
		return nil, err
	}

	var pool IPPool
	if err := c.Do(ctx, req, &pool); err != nil {
		return nil, err
	}

	return &pool, nil
}

// DeleteIPPool deletes an IP pool
// see: https://www.twilio.com/docs/sendgrid/api-reference/ip-pools/delete-an-ip-pool
func (c *Client) DeleteIPPool(ctx context.Context, name string) error {
	path := fmt.Sprintf("/ips/pools/%s", url.QueryEscape(name))

	req, err := c.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	if err := c.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// AddIPToPool adds an IP address to a pool
// see: https://www.twilio.com/docs/sendgrid/api-reference/ip-pools/add-an-ip-address-to-a-pool
func (c *Client) AddIPToPool(ctx context.Context, poolName, ip string) error {
	path := fmt.Sprintf("/ips/pools/%s/ips", url.QueryEscape(poolName))

	input := &InputAddIPToPool{IP: ip}

	req, err := c.NewRequest("POST", path, input)
	if err != nil {
		return err
	}

	if err := c.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// RemoveIPFromPool removes an IP address from a pool
// see: https://www.twilio.com/docs/sendgrid/api-reference/ip-pools/remove-an-ip-address-from-a-pool
func (c *Client) RemoveIPFromPool(ctx context.Context, poolName, ip string) error {
	path := fmt.Sprintf("/ips/pools/%s/ips/%s", url.QueryEscape(poolName), url.QueryEscape(ip))

	req, err := c.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	if err := c.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// StartIPWarmup starts IP warmup process
// see: https://www.twilio.com/docs/sendgrid/api-reference/ip-warmup/start-ip-warmup
func (c *Client) StartIPWarmup(ctx context.Context, ip string) (*IPWarmupStatus, error) {
	path := fmt.Sprintf("/ips/warmup/%s", url.QueryEscape(ip))

	req, err := c.NewRequest("POST", path, nil)
	if err != nil {
		return nil, err
	}

	var status IPWarmupStatus
	if err := c.Do(ctx, req, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

// StopIPWarmup stops IP warmup process
// see: https://www.twilio.com/docs/sendgrid/api-reference/ip-warmup/stop-ip-warmup
func (c *Client) StopIPWarmup(ctx context.Context, ip string) (*IPWarmupStatus, error) {
	path := fmt.Sprintf("/ips/warmup/%s", url.QueryEscape(ip))

	req, err := c.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	var status IPWarmupStatus
	if err := c.Do(ctx, req, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

// GetIPWarmupStatus retrieves IP warmup status
// see: https://www.twilio.com/docs/sendgrid/api-reference/ip-warmup/retrieve-ip-warmup-status
func (c *Client) GetIPWarmupStatus(ctx context.Context, ip string) (*IPWarmupStatus, error) {
	path := fmt.Sprintf("/ips/warmup/%s", url.QueryEscape(ip))

	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var status IPWarmupStatus
	if err := c.Do(ctx, req, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

// GetAllIPWarmupStatus retrieves all IP warmup statuses
// see: https://www.twilio.com/docs/sendgrid/api-reference/ip-warmup/retrieve-all-ip-warmup-statuses
func (c *Client) GetAllIPWarmupStatus(ctx context.Context) ([]IPWarmupStatus, error) {
	path := "/ips/warmup"

	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var statuses []IPWarmupStatus
	if err := c.Do(ctx, req, &statuses); err != nil {
		return nil, err
	}

	return statuses, nil
}
