package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jobstoit/hetzner-dns-go/dns/schema"
)

func TestZoneList(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc(pathZones, func(w http.ResponseWriter, r *http.Request) {
		res := schema.ZoneListResponse{}
		res.Zones = []schema.Zone{
			{
				ID:   "1",
				Name: "hetzner.com",
			},
			{
				ID:   "2",
				Name: "hetzner.cloud",
			},
		}

		if r.URL.Query().Get("name") != "" {
			res.Zones = []schema.Zone{
				res.Zones[0],
			}
		}

		json.NewEncoder(w).Encode(res) // nolint: errcheck
	}).Methods("GET")

	ctx := context.Background()
	opts := ZoneListOpts{}
	zones, _, err := env.Client.Zone.List(ctx, opts)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(zones) != 2 {
		t.Errorf("expected %d but got %d", 2, len(zones))
	}

	opts.Name = "hetzner.com"
	zones, _, err = env.Client.Zone.List(ctx, opts)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(zones) != 1 {
		t.Errorf("expected %d but got %d", 1, len(zones))
	}
}

func TestZoneGetByID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc(fmt.Sprintf("%s/{id}", pathZones), func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		paramID := params["id"]

		if paramID != "1" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		resp := schema.Zone{
			ID: paramID,
		}

		json.NewEncoder(w).Encode(schema.ZoneResponse{ // nolint: errcheck
			Zone: resp,
		})
	}).Methods("GET")

	ctx := context.Background()

	id := "1"
	zone, _, err := env.Client.Zone.GetByID(ctx, id)
	if err != nil {
		t.Errorf("error fetching server: %v", err)
	}

	if zone == nil || zone.ID != id {
		t.Error("missing zone")
	}

	id = "0"
	_, resp, err := env.Client.Zone.GetByID(ctx, id)
	if err == nil {
		t.Error("missing expected error")
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("unexpected status code, expected %d but got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestZoneCreate(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	zoneID := "2"

	env.Mux.HandleFunc(pathZones, func(w http.ResponseWriter, r *http.Request) {
		var body schema.ZoneCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TestWrongValue
		if body.Name == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		var res schema.ZoneResponse
		res.Zone = schema.Zone{
			ID:   zoneID,
			Name: body.Name,
			Ttl:  *body.Ttl,
		}

		json.NewEncoder(w).Encode(res) // nolint: errcheck
	}).Methods("POST")

	ctx := context.Background()

	opts := ZoneCreateOpts{
		Name: "",
		Ttl:  nil,
	}

	_, _, err := env.Client.Zone.Create(ctx, opts)
	if err == nil {
		t.Error("error expected but got nil")
	}

	opts.Name = "hetzner.com"
	ttl := 86400
	opts.Ttl = &ttl

	zone, _, err := env.Client.Zone.Create(ctx, opts)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if zone.ID != zoneID {
		t.Errorf("expected '%s' but got '%s'", zoneID, zone.ID)
	}
}

func TestZoneUpdate(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc(fmt.Sprintf("%s/{id}", pathZones), func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		paramID := params["id"]

		if paramID != "1" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var body schema.ZoneUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var resp schema.ZoneResponse
		resp.Zone = schema.Zone{
			ID:   paramID,
			Name: body.Name,
		}

		json.NewEncoder(w).Encode(resp) // nolint: errcheck
	}).Methods("PUT")

	ctx := context.Background()
	zone := &Zone{
		ID: "0",
	}
	opts := ZoneUpdateOpts{
		Name: "",
	}

	_, _, err := env.Client.Zone.Update(ctx, zone, opts)
	if err.Error() != "name required" {
		t.Errorf("expected error 'name required' but got '%v'", err)
	}

	opts.Name = "hetzner.com"
	_, resp, err := env.Client.Zone.Update(ctx, zone, opts)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected %d but got %d", http.StatusNotFound, resp.StatusCode)
	} else if resp.StatusCode != http.StatusNotFound && err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	zone.ID = "1"

	updatedZone, _, err := env.Client.Zone.Update(ctx, zone, opts)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if updatedZone.ID != zone.ID {
		t.Errorf("expected '%s' but got '%s'", zone.ID, updatedZone.ID)
	}
}

func TestZoneDelete(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc(fmt.Sprintf("%s/{id}", pathZones), func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		paramID := params["id"]

		if paramID != "1" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	})

	ctx := context.Background()
	zone := &Zone{
		ID: "0",
	}

	resp, err := env.Client.Zone.Delete(ctx, zone)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected %d but got %d", http.StatusNotFound, resp.StatusCode)
	} else if resp.StatusCode != http.StatusNotFound && err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	zone.ID = "1"
	_, err = env.Client.Zone.Delete(ctx, zone)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TODO create tests for Import and Export and validate
