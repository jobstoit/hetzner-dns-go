package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/jobstoit/hetzner-dns-go/dns/schema"
)

// PrimaryServer represents a primary server in the Hetzner DNS API.
type PrimaryServer struct {
	ID       string
	Port     int
	Created  schema.HdnsTime
	Modified schema.HdnsTime
	Zone     *Zone
	Address  string
}

// RecordClient is a client for primary server API.
type PrimaryServerClient struct {
	client *Client
}

// PrimaryServerListOpts specifies options for listing primary servers
type PrimaryServerListOpts struct {
	ZoneID string
}

func (o PrimaryServerListOpts) values() url.Values {
	vals := url.Values{}
	if o.ZoneID != "" {
		vals.Add("zone_id", o.ZoneID)
	}
	return vals
}

// List returns all primary servers with the given parameters.
func (c PrimaryServerClient) List(ctx context.Context, opts PrimaryServerListOpts) ([]*PrimaryServer, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("%s?%s", pathPrimaryServers, opts.values().Encode()), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.PrimaryServerListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	primaryServers := make([]*PrimaryServer, 0, len(body.PrimaryServers))
	for _, server := range body.PrimaryServers {
		primaryServers = append(primaryServers, PrimaryServerFromSchema(server))
	}

	return primaryServers, resp, nil
}

// GetByID returns the PrimaryServer with the given id.
func (c PrimaryServerClient) GetByID(ctx context.Context, id string) (*PrimaryServer, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("%s/%s", pathPrimaryServers, id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.PrimaryServerResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	return PrimaryServerFromSchema(body.PrimaryServer), resp, nil
}

// PrimaryServerCreateOpts specifies options for creating a primary server.
type PrimaryServerCreateOpts struct {
	Address string
	Port    int
	ZoneID  string
}

func (o PrimaryServerCreateOpts) validate() error {
	if o.Address == "" {
		return errors.New("address required")
	}
	if o.Port < 1 || o.Port > 65535 {
		return errors.New("invalid port")
	}
	if o.ZoneID == "" {
		return errors.New("zone_id required")
	}
	return nil
}

// Create creates a new primary server record.
func (c PrimaryServerClient) Create(ctx context.Context, opts PrimaryServerCreateOpts) (*PrimaryServer, *Response, error) {
	if err := opts.validate(); err != nil {
		return nil, nil, err
	}

	var reqBody schema.PrimaryServerCreateRequest
	reqBody.Address = opts.Address
	reqBody.Port = opts.Port
	reqBody.ZoneID = opts.ZoneID

	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.client.NewRequest(ctx, "POST", pathPrimaryServers, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var body schema.PrimaryServerResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}

	return PrimaryServerFromSchema(body.PrimaryServer), resp, nil
}

// PrimaryServerUpdateOpts specifies options for updating a primary server.
type PrimaryServerUpdateOpts struct {
	Address string
	Port    int
	ZoneID  string
}

func (o PrimaryServerUpdateOpts) validate() error {
	if o.Address == "" {
		return errors.New("address required")
	}
	if o.Port < 1 || o.Port > 65535 {
		return errors.New("invalid port")
	}
	if o.ZoneID == "" {
		return errors.New("zone_id required")
	}
	return nil
}

// Update updates a primary server record.
func (c PrimaryServerClient) Update(ctx context.Context, server *PrimaryServer, opts PrimaryServerUpdateOpts) (*PrimaryServer, *Response, error) {
	if err := opts.validate(); err != nil {
		return nil, nil, err
	}

	var reqBody schema.PrimaryServerUpdateRequest
	reqBody.Address = opts.Address
	reqBody.Port = opts.Port
	reqBody.ZoneID = opts.ZoneID

	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.client.NewRequest(ctx, "PUT", fmt.Sprintf("%s/%s", pathPrimaryServers, server.ID), bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var body schema.PrimaryServerResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	return PrimaryServerFromSchema(body.PrimaryServer), resp, nil
}

// Delete deletes a primary server record.
func (c PrimaryServerClient) Delete(ctx context.Context, server *PrimaryServer) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("%s/%s", pathPrimaryServers, server.ID), nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req, nil)
}
