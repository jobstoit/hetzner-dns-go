// Package dns an SDK for the Hetzner DNS API.
//
// Read the API docs over on https://dns.hetzner.com/api-docs.
//
// Example:
//
//   package main
//
//   import (
//   	"context"
//   	"fmt"
//   	"log"
//
//   	"github.com/jobstoit/hetzner-dns-go/dns"
//   )
//
//   func main() {
//   	client := dns.NewClient(dns.WithToken("token"))
//
//   	record, _, err := client.Record.GetByID(context.Background(), "randomid")
//   	if err != nil {
//   		log.Fatalf("error retrieving record: %v\n", err)
//   	}
//
//   	fmt.Printf("record of type: '%s' found with value: %s", record.Type, record.Value)
//   }
//
package dns
