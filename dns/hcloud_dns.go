package dns

// Version of the SDK
const Version = "v0.1.0"

// Endpoint is the base URL of the API.
const Endpoint = "https://dns.hetzner.com/api/v1"

const UserAgent = "hetzner-dns/" + Version

const (
	pathZones          = "/zones"
	pathRecords        = "/records"
	pathPrimaryServers = "/primary_servers"
)
