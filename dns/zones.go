package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/jobstoit/hetzner-dns-go/dns/schema"
)

type ZoneStatus string

const (
	ZoneStatusVerified ZoneStatus = `verified`
	ZoneStatusFailed   ZoneStatus = `failed`
	ZoneStatusPending  ZoneStatus = `pending`
)

// Zone represents a zone in Hetzner DNS.
type Zone struct {
	ID              string
	Created         schema.HdnsTime
	Modified        schema.HdnsTime
	LegacyDNSHost   string
	LegacyNS        []string
	Name            string
	NS              []string
	Owner           string
	Paused          bool
	Permission      string
	Project         string
	Registrar       string
	Status          ZoneStatus
	Ttl             int
	Verified        schema.HdnsTime
	RecordsCount    int
	IsSecondaryDNS  bool
	TxtVerification *TxtVerification
}

// TxtVerification represents a txt verification of a zone.
type TxtVerification struct {
	Name  string
	Token string
}

// ZoneClient is a client for zones API.
type ZoneClient struct {
	client *Client
}

//  ZoneListOptions specifies options for listing zones.
type ZoneListOpts struct {
	ListOpts
	Name       string
	SearchName string
}

func (l ZoneListOpts) values() url.Values {
	vals := l.ListOpts.values()
	if l.Name != "" {
		vals.Add("name", l.Name)
	}
	if l.SearchName != "" {
		vals.Add("search_name", l.SearchName)
	}
	return vals
}

// List returns all zones with the given parameters.
func (c ZoneClient) List(ctx context.Context, opts ZoneListOpts) ([]*Zone, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("%s?%s", pathZones, opts.values().Encode()), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ZoneListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	zones := make([]*Zone, 0, len(body.Zones))
	for _, z := range body.Zones {
		zones = append(zones, ZoneFromSchema(z))
	}

	return zones, resp, nil
}

// GetByID returns the zone with the given id.
func (c ZoneClient) GetByID(ctx context.Context, id string) (*Zone, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("%s/%s", pathZones, id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ZoneResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	return ZoneFromSchema(body.Zone), resp, nil
}

// ZoneCreateOpts specifies options for creating a new zone.
type ZoneCreateOpts struct {
	Name string
	Ttl  *int
}

// Validate checks if the options are valid.
func (o ZoneCreateOpts) Validate() error {
	if o.Name == "" {
		return errors.New("name required")
	}

	return nil
}

// Create creates a new zone.
func (c ZoneClient) Create(ctx context.Context, opts ZoneCreateOpts) (*Zone, *Response, error) {
	if err := opts.Validate(); err != nil {
		return nil, nil, err
	}

	var reqBody schema.ZoneCreateRequest
	reqBody.Name = opts.Name
	reqBody.Ttl = opts.Ttl

	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.client.NewRequest(ctx, "POST", pathZones, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var body schema.ZoneResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	return ZoneFromSchema(body.Zone), resp, nil
}

// ZoneUpdateOpts specifies options for updating a zone.
type ZoneUpdateOpts struct {
	Name string
	Ttl  *int
}

// Validate checks if the options are valid.
func (o ZoneUpdateOpts) Validate() error {
	if o.Name == "" {
		return errors.New("name required")
	}

	return nil
}

// Update updates a zone.
func (c ZoneClient) Update(ctx context.Context, zone *Zone, opts ZoneUpdateOpts) (*Zone, *Response, error) {
	if err := opts.Validate(); err != nil {
		return nil, nil, err
	}

	var reqBody schema.ZoneUpdateRequest
	reqBody.Name = opts.Name
	reqBody.Ttl = opts.Ttl

	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.client.NewRequest(ctx, "PUT", fmt.Sprintf("%s/%s", pathZones, zone.ID), bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var body schema.ZoneResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	return ZoneFromSchema(body.Zone), resp, nil
}

// Delete deletes a zone.
func (c ZoneClient) Delete(ctx context.Context, zone *Zone) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("%s/%s", pathZones, zone.ID), nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req, nil)
}

// Import imports a zone file in text/plain format.
func (c ZoneClient) Import(ctx context.Context, zone *Zone, file io.Reader) (*Zone, *Response, error) {
	req, err := c.client.NewRequest(ctx, "POST", fmt.Sprintf("%s/%s/import", pathZones, zone.ID), file)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "text/plain")

	var body schema.ZoneResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	return ZoneFromSchema(body.Zone), resp, nil
}

// Export exports a zone in text/plain format.
func (c ZoneClient) Export(ctx context.Context, zone *Zone) (io.Reader, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("%s/%s/export", pathZones, zone.ID), nil)
	if err != nil {
		return nil, nil, err
	}

	file := &bytes.Buffer{}
	resp, err := c.client.Do(req, file)
	if err != nil {
		return nil, resp, err
	}

	return file, resp, nil
}

// ValidatedZoneFile is returned when validating a zone file.
type ValidatedZoneFile struct {
	PassedRecords int
	ValidRecords  []*Record
}

// ValidateFile validates a zone file.
func (c ZoneClient) ValidateFile(ctx context.Context, file io.Reader) (*ValidatedZoneFile, *Response, error) {
	req, err := c.client.NewRequest(ctx, "POST", fmt.Sprintf("%s/file/validate", pathZones), file)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "text/plain")

	var body schema.ValidateZoneFileResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	return ValidateZoneFileFromSchema(body), resp, nil
}
