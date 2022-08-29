package schema

import "time"

// Record represents a record in Hetzner DNS.
type Record struct {
	Type     string    `json:"type"`
	ID       string    `json:"id"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	ZoneID   string    `json:"zone_id"`
	Name     string    `json:"name"`
	Value    string    `json:"value"`
	Ttl      int       `json:"ttl"`
}

// RecordListResponse defines the schema of the response when
// listing zones.
type RecordListResponse struct {
	Records []Record `json:"records"`
}

// RecordResponse defines the schema of the response when
// listing zones.
type RecordResponse struct {
	Record Record `json:"record"`
}

// RecordCreateRequest defines a schema for the request to
// create a record.
type RecordCreateRequest struct {
	Name   string `json:"name"`
	Ttl    *int   `json:"ttl"`
	Type   string `json:"type"`
	Value  string `json:"value"`
	ZoneID string `json:"zone_id"`
}

// RecordUpdateRequest defines a schema for the request to
// update a record.
type RecordUpdateRequest struct {
	Name   string `json:"name"`
	Ttl    *int   `json:"ttl"`
	Type   string `json:"type"`
	Value  string `json:"value"`
	ZoneID string `json:"zone_id"`
}

// RecordBulkCreateRequest defines a schema for the request to
// bulk create records.
type RecordBulkCreateRequest struct {
	Records []RecordCreateRequest `json:"records"`
}

// RecordResponse defines the schema of the response when
// listing zones.
type RecordBulkCreateResponse struct {
	Records        []Record          `json:"record"`
	ValidRecords   []RecordBulkEntry `json:"valid_records"`
	InvalidRecords []RecordBulkEntry `json:"invalid_records"`
}

// RecordBulkEntry defines a schema for an entry for the request to
// bulk create records and is used in responses.
type RecordBulkEntry struct {
	Name   string `json:"name"`
	Ttl    *int   `json:"ttl"`
	Type   string `json:"type"`
	Value  string `json:"value"`
	ZoneID string `json:"zone_id"`
}

// RecordBulkUpdateRequest defines a schema for the request to
// bulk update records.
type RecordBulkUpdateRequest struct {
	Records []RecordBulkUpdateEntry `json:"records"`
}

// RecordBulkUpdateEntry defines a schema for an entry for the request to
// bulk update records.
type RecordBulkUpdateEntry struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Ttl    *int   `json:"ttl"`
	Type   string `json:"type"`
	Value  string `json:"value"`
	ZoneID string `json:"zone_id"`
}

// RecordBulkUpdateResponse defines a schema for the respose to
// bulk update records.
type RecordBulkUpdateResponse struct {
	Records       []Record          `json:"records"`
	FailedRecords []RecordBulkEntry `json:"failed_records"`
}
