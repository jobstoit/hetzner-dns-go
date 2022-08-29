package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/jobstoit/hcloud-dns-go/dns/schema"
)

type RecordType string

const (
	RecordTypeA     = "A"
	RecordTypeAAAA  = "AAAA"
	RecordTypePTR   = "PTR"
	RecordTypeNS    = "NS"
	RecordTypeMX    = "MX"
	RecordTypeCNAME = "CNAME"
	RecordTypeRP    = "RP"
	RecordTypeTXT   = "TXT"
	RecordTypeSOA   = "SOA"
	RecordTypeHINFO = "HINFO"
	RecordTypeSRV   = "SRV"
	RecordTypeDANE  = "DANE"
	RecordTypeTLSA  = "TLSA"
	RecordTypeDS    = "DS"
	RecordTypeCAA   = "CAA"
)

// Record represents a record in the Hetzner DNS.
type Record struct {
	Type     RecordType
	ID       string
	Created  time.Time
	Modified time.Time
	Zone     *Zone
	Name     string
	Value    string
	Ttl      int
}

// RecordClient is a client for records API.
type RecordClient struct {
	client *Client
}

// RecordListOpts specifies options for listing records
type RecordListOpts struct {
	ListOpts
	ZoneID string
}

func (o RecordListOpts) values() url.Values {
	vals := o.ListOpts.values()
	if o.ZoneID != "" {
		vals.Add("zone_id", o.ZoneID)
	}
	return vals
}

// List returns all records with the given parameters.
func (c RecordClient) List(ctx context.Context, opts RecordListOpts) ([]*Record, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("%s?%s", pathRecords, opts.values().Encode()), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.RecordListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}

	records := make([]*Record, 0, len(body.Records))
	for _, r := range body.Records {
		records = append(records, RecordFromSchema(r))
	}

	return records, resp, nil
}

// GetByID returns a record with the given id.
func (c RecordClient) GetByID(ctx context.Context, id string) (*Record, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("%s/%s", pathRecords, id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.RecordResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}

	return RecordFromSchema(body.Record), resp, nil
}

// RecordCreateOpts specifies options for creating a record.
type RecordCreateOpts struct {
	Name  string
	Ttl   *int
	Type  RecordType
	Value string
	Zone  *Zone
}

func (o RecordCreateOpts) validate() error {
	if o.Name == "" {
		return errors.New("name required")
	}
	if o.Type == "" {
		return errors.New("type required")
	}
	if o.Value == "" {
		return errors.New("value required")
	}
	if o.Zone == nil {
		return errors.New("zone required")
	}

	return nil
}

// Create creates a new record.
func (c RecordClient) Create(ctx context.Context, opts RecordCreateOpts) (*Record, *Response, error) {
	if err := opts.validate(); err != nil {
		return nil, nil, err
	}

	var reqBody schema.RecordCreateRequest
	reqBody.Name = opts.Name
	reqBody.Ttl = opts.Ttl
	reqBody.Type = string(opts.Type)
	reqBody.Value = opts.Value
	reqBody.ZoneID = opts.Zone.ID

	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.client.NewRequest(ctx, "POST", pathRecords, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var body schema.RecordResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	return RecordFromSchema(body.Record), resp, nil
}

// RecordUpdateOpts specifies options for creating a record.
type RecordUpdateOpts struct {
	Name  string
	Ttl   *int
	Type  RecordType
	Value string
	Zone  *Zone
}

func (o RecordUpdateOpts) validate() error {
	if o.Name == "" {
		return errors.New("name required")
	}
	if o.Type == "" {
		return errors.New("type required")
	}
	if o.Value == "" {
		return errors.New("value required")
	}
	if o.Zone == nil {
		return errors.New("zone required")
	}

	return nil
}

// Update updates a record.
func (c RecordClient) Update(ctx context.Context, rec *Record, opts RecordUpdateOpts) (*Record, *Response, error) {
	if err := opts.validate(); err != nil {
		return nil, nil, err
	}

	var reqBody schema.RecordUpdateRequest
	reqBody.Name = opts.Name
	reqBody.Ttl = opts.Ttl
	reqBody.Type = string(opts.Type)
	reqBody.Value = opts.Value
	reqBody.ZoneID = rec.Zone.ID

	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.client.NewRequest(ctx, "PUT", fmt.Sprintf("%s/%s", pathRecords, rec.ID), bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var body schema.RecordResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	return RecordFromSchema(body.Record), resp, nil
}

// Delete deletes a record.
func (c RecordClient) Delete(ctx context.Context, rec *Record) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("%s/%s", pathRecords, rec.ID), nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req, nil)
}

// RecordEntry represents a record entry used for verification and updates.
type RecordEntry struct {
	Type   RecordType
	ZoneID string
	Name   string
	Value  string
	Ttl    *int
}

// RecordBulkCreateResponse is returned when creating records in bulk.
type RecordBulkCreateResponse struct {
	Records        []*Record
	ValidRecords   []*RecordEntry
	InvalidRecords []*RecordEntry
}

// BulkCreate creates multiple records.
func (c RecordClient) BulkCreate(ctx context.Context, bulkOpts []RecordCreateOpts) (*RecordBulkCreateResponse, *Response, error) {
	for _, opt := range bulkOpts {
		if err := opt.validate(); err != nil {
			return nil, nil, err
		}
	}

	var reqBody schema.RecordBulkCreateRequest
	reqBody.Records = make([]schema.RecordCreateRequest, 0, len(bulkOpts))
	for _, opt := range bulkOpts {
		var r schema.RecordCreateRequest
		r.Name = opt.Name
		r.Ttl = opt.Ttl
		r.Type = string(opt.Type)
		r.Value = opt.Value
		r.ZoneID = opt.Zone.ID

		reqBody.Records = append(reqBody.Records, r)
	}

	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.client.NewRequest(ctx, "POST", fmt.Sprintf("%s/bulk", pathRecords), bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var body schema.RecordBulkCreateResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	respBody := &RecordBulkCreateResponse{}
	for _, rec := range body.Records {
		respBody.Records = append(respBody.Records, RecordFromSchema(rec))
	}

	for _, rec := range body.ValidRecords {
		respBody.ValidRecords = append(respBody.ValidRecords, RecordEntryFromSchema(rec))
	}

	for _, rec := range body.InvalidRecords {
		respBody.InvalidRecords = append(respBody.InvalidRecords, RecordEntryFromSchema(rec))
	}

	return respBody, resp, nil
}

// RecordBulkUpdateOpts specifies options for a bulk record update entry
type RecordBulkUpdateOpts struct {
	ID    string
	Type  RecordType
	Zone  *Zone
	Name  string
	Value string
	Ttl   *int
}

func (o RecordBulkUpdateOpts) validate() error {
	if o.ID == "" {
		return errors.New("id required")
	}
	if o.Name == "" {
		return errors.New("name required")
	}
	if o.Type == "" {
		return errors.New("type required")
	}
	if o.Value == "" {
		return errors.New("value required")
	}
	if o.Zone == nil {
		return errors.New("zone required")
	}

	return nil
}

// RecordBulkUpdateResponse is returned when creating records in bulk.
type RecordBulkUpdateResponse struct {
	Records       []*Record
	FailedRecords []*RecordEntry
}

// BulkUpdate updates multiple records.
func (c RecordClient) BulkUpdate(ctx context.Context, bulkOpts []RecordBulkUpdateOpts) (*RecordBulkUpdateResponse, *Response, error) {
	for _, opts := range bulkOpts {
		if err := opts.validate(); err != nil {
			return nil, nil, err
		}
	}

	var reqBody schema.RecordBulkUpdateRequest
	reqBody.Records = make([]schema.RecordBulkUpdateEntry, 0, len(bulkOpts))
	for _, opts := range bulkOpts {
		var recBody schema.RecordBulkUpdateEntry
		recBody.ID = opts.ID
		recBody.Name = opts.Name
		recBody.Type = string(opts.Type)
		recBody.Value = opts.Value
		recBody.Ttl = opts.Ttl
		recBody.ZoneID = opts.Zone.ID

		reqBody.Records = append(reqBody.Records, recBody)
	}

	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.client.NewRequest(ctx, "PUT", fmt.Sprintf("%s/bulk", pathRecords), bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var body schema.RecordBulkUpdateResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	respBody := &RecordBulkUpdateResponse{}
	for _, rec := range body.Records {
		respBody.Records = append(respBody.Records, RecordFromSchema(rec))
	}

	for _, rec := range body.FailedRecords {
		respBody.FailedRecords = append(respBody.FailedRecords, RecordEntryFromSchema(rec))
	}

	return respBody, resp, nil
}
