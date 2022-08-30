package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

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
	})

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

	as := newAssert(t)

	env.Mux.HandleFunc(fmt.Sprintf("%s/1", pathZones), func(w http.ResponseWriter, r *http.Request) {
		resp := schema.Zone{
			ID: "1",
		}

		json.NewEncoder(w).Encode(schema.ZoneResponse{ // nolint: errcheck
			Zone: resp,
		})
	})

	ctx := context.Background()

	id := "0"
	_, resp, err := env.Client.Zone.GetByID(ctx, id)
	if as.EqInt(http.StatusNotFound, resp.StatusCode) {
		as.Error(err)
	}

	id = "1"
	zone, _, err := env.Client.Zone.GetByID(ctx, id)
	if as.NoError(err) && as.NotNil(zone) {
		as.EqStr(id, zone.ID)
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

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res) // nolint: errcheck
	})

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

	as := newAssert(t)

	env.Mux.HandleFunc(fmt.Sprintf("%s/1", pathZones), func(w http.ResponseWriter, r *http.Request) {
		var body schema.ZoneUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var resp schema.ZoneResponse
		resp.Zone = schema.Zone{
			ID:   "1",
			Name: body.Name,
		}

		json.NewEncoder(w).Encode(resp) // nolint: errcheck
	})

	ctx := context.Background()
	zone := &Zone{
		ID: "0",
	}
	opts := ZoneUpdateOpts{
		Name: "",
	}

	_, _, err := env.Client.Zone.Update(ctx, zone, opts)
	as.EqStr("name required", err.Error())

	opts.Name = "hetzner.com"
	_, resp, err := env.Client.Zone.Update(ctx, zone, opts)
	if !as.EqInt(http.StatusNotFound, resp.StatusCode) {
		as.NoError(err)
	}

	zone.ID = "1"

	updatedZone, _, err := env.Client.Zone.Update(ctx, zone, opts)
	as.NoError(err)
	as.EqStr(zone.ID, updatedZone.ID)
	as.EqStr(opts.Name, updatedZone.Name)

}

func TestZoneDelete(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	as := newAssert(t)

	env.Mux.HandleFunc(fmt.Sprintf("%s/1", pathZones), func(w http.ResponseWriter, r *http.Request) {})

	ctx := context.Background()
	zone := &Zone{
		ID: "0",
	}

	resp, err := env.Client.Zone.Delete(ctx, zone)
	if !as.EqInt(http.StatusNotFound, resp.StatusCode) {
		as.NoError(err)
	}

	zone.ID = "1"
	_, err = env.Client.Zone.Delete(ctx, zone)
	as.NoError(err)
}

func TestZoneImport(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	as := newAssert(t)

	env.Mux.HandleFunc(fmt.Sprintf("%s/1/import", pathZones), func(w http.ResponseWriter, r *http.Request) {
		body := &bytes.Buffer{}
		body.ReadFrom(r.Body) // nolint: errcheck

		if body.String() != "valid syntax" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		var respBody schema.ZoneResponse
		respBody.Zone.ID = "1"

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(respBody) // nolint: errcheck
	})

	ctx := context.Background()
	file := bytes.NewBufferString("invalid syntax")
	zone := &Zone{ID: "0"}

	_, resp, err := env.Client.Zone.Import(ctx, zone, file)
	if !as.EqInt(http.StatusNotFound, resp.StatusCode) {
		as.NoError(err)
	}

	zone.ID = "1"
	_, resp, err = env.Client.Zone.Import(ctx, zone, file)
	if !as.EqInt(http.StatusUnprocessableEntity, resp.StatusCode) {
		as.NoError(err)
	}

	file = bytes.NewBufferString("valid syntax")
	newZone, _, err := env.Client.Zone.Import(ctx, zone, file)
	if as.NoError(err) {
		as.EqStr(zone.ID, newZone.ID)
	}
}

func TestZoneExport(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	as := newAssert(t)

	env.Mux.HandleFunc(fmt.Sprintf("%s/1/export", pathZones), func(w http.ResponseWriter, r *http.Request) {
		respBody := bytes.NewBufferString("valid syntax")

		respBody.WriteTo(w) // nolint: errcheck
	})

	ctx := context.Background()
	zone := &Zone{ID: "0"}

	_, resp, err := env.Client.Zone.Export(ctx, zone)
	if !as.EqInt(http.StatusNotFound, resp.StatusCode) {
		as.NoError(err)
	}

	zone.ID = "1"
	file, _, err := env.Client.Zone.Export(ctx, zone)
	if as.NoError(err) {
		body := &bytes.Buffer{}
		body.ReadFrom(file) // nolint: errcheck

		as.EqStr("valid syntax", body.String())
	}
}

func TestZoneValidate(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	as := newAssert(t)

	env.Mux.HandleFunc(fmt.Sprintf("%s/file/validate", pathZones), func(w http.ResponseWriter, r *http.Request) {
		body := &bytes.Buffer{}
		body.ReadFrom(r.Body) // nolint: errcheck

		if body.String() != "valid syntax" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		var respBody schema.ValidateZoneFileResponse
		respBody.PassedRecords = 2
		respBody.ValidRecords = []schema.Record{
			{
				Name: "hetzner.com",
			},
			{
				Name: "dns.hetzner.com",
			},
		}

		json.NewEncoder(w).Encode(respBody) // nolint: errcheck
	})

	ctx := context.Background()
	file := bytes.NewBufferString("invalid syntax")
	_, resp, err := env.Client.Zone.ValidateFile(ctx, file)
	if !as.EqInt(http.StatusUnprocessableEntity, resp.StatusCode) {
		as.NoError(err)
	}

	file = bytes.NewBufferString("valid syntax")
	val, _, err := env.Client.Zone.ValidateFile(ctx, file)
	if as.NoError(err) {
		as.EqInt(2, val.PassedRecords)
		as.EqInt(2, len(val.ValidRecords))
	}
}
