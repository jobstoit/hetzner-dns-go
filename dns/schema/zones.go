package schema

// ZoneListResponse defines the schema of the response when
// listing zones.
type ZoneListResponse struct {
	Zones []Zone `json:"zones"`
}

// ZoneGetResponse defines the schema of the response when
// listing zones.
type ZoneResponse struct {
	Zone Zone `json:"zone"`
}

// Zone represents a zone in Hetzner DNS.
type Zone struct {
	ID              string          `json:"id"`
	Created         HdnsTime        `json:"created"`
	Modified        HdnsTime        `json:"modified"`
	LegacyDNSHost   string          `json:"legacy_dns_host"`
	LegacyNS        []string        `json:"legacy_ns"`
	Name            string          `json:"name"`
	NS              []string        `json:"ns"`
	Owner           string          `json:"owner"`
	Paused          bool            `json:"paused"`
	Permission      string          `json:"permission"`
	Project         string          `json:"project"`
	Registrar       string          `json:"registrar"`
	Status          string          `json:"status"`
	Ttl             int             `json:"ttl"`
	Verified        HdnsTime        `json:"verified"`
	RecordsCount    int             `json:"records_count"`
	IsSecondaryDNS  bool            `json:"is_secondary_dns"`
	TxtVerification TxtVerification `json:"txt_verification"`
}

// TxtVerification represents the text verification of a zone.
type TxtVerification struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

// CreateZoneRequest defines a schema for the request to
// create a zone.
type ZoneCreateRequest struct {
	Name string `json:"name"`
	Ttl  *int   `json:"ttl"`
}

// CreateZoneRequest defines a schema for the request to
// create a zone.
type ZoneUpdateRequest struct {
	Name string `json:"name"`
	Ttl  *int   `json:"ttl"`
}

// ValidateZoneFileResponse defines the schema of the response when
// validating a zone file.
type ValidateZoneFileResponse struct {
	PassedRecords int      `json:"passed_records"`
	ValidRecords  []Record `json:"valid_records"`
}
