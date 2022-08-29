# Hetzner DNS: a Go library for the Hetzner DNS API
[![Go Reference](https://pkg.go.dev/badge/github.com/jobstoit/hetzner-dns-go/dns.svg)](https://pkg.go.dev/github.com/jobstoit/hetzner-dns-go/dns)

Package dns is a library for the Hetzner DNS API.

The libraries documentation is available at [Go Pkg](https://pkg.go.dev/github.com/jobstoit/hetzner-dns-go/dns), the public API documentation is available at [dns.hetzner.com](https://dns.hetzner.com/api-docs).

## Example
```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jobstoit/hetzner-dns-go/dns"
)

func main() {
	client := dns.NewClient(dns.WithToken("token"))

	record, _, err := client.Record.GetByID(context.Background(), "randomid")
	if err != nil {
		log.Fatalf("error retrieving record: %v\n", err)
	}

	fmt.Printf("record of type: '%s' found with value: %s", record.Type, record.Value)
}
```
