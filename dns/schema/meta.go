package schema

// MetaResponse defines the schema of a response containing
// meta information.
type MetaResponse struct {
	Meta Meta `json:"meta"`
}

// Meta defines the schema of meta information which may be included
// in responses.
type Meta struct {
	Pagination *MetaPagination `json:"pagination"`
}

// MetaPagination defines the schema of pagination information.
type MetaPagination struct {
	LastPage     int `json:"last_page"`
	Page         int `json:"page"`
	PerPage      int `json:"per_page"`
	TotalEntries int `json:"total_entries"`
}
