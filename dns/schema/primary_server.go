package schema

import "time"

// PrimaryServer represents a primary server in the Hetzner DNS API.
type PrimaryServer struct {
	ID       string    `json:"id"`
	Port     int       `json:"port"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	ZoneID   string    `json:"zone_id"`
	Address  string    `json:"address"`
}

// PrimaryServerListResponse defines the schema of the response when
// listing primary servers.
type PrimaryServerListResponse struct {
	PrimaryServers []PrimaryServer `json:"primary_servers"`
}

// PrimaryServerResponse defines the schema of the response when
// listing zones.
type PrimaryServerResponse struct {
	PrimaryServer PrimaryServer `json:"primary_server"`
}

// PrimaryServerCreateRequest defines a schema for the request to
// create a primary server.
type PrimaryServerCreateRequest struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	ZoneID  string `json:"zone_id"`
}

// PrimaryServerUpdateRequest defines a schema for the request to
// create a primary server.
type PrimaryServerUpdateRequest struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	ZoneID  string `json:"zone_id"`
}
