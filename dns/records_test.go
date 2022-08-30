package dns

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jobstoit/hetzner-dns-go/dns/schema"
)

func TestRecordList(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	as := newAssert(t)

	env.Mux.HandleFunc(pathRecords, func(w http.ResponseWriter, r *http.Request) {
		resp := schema.RecordListResponse{
			Records: []schema.Record{
				{
					ID:     "1",
					Type:   "A",
					Name:   "hetzner.cloud",
					ZoneID: "1",
				},
				{
					ID:     "2",
					Type:   "A",
					Name:   "hetzner.com",
					ZoneID: "2",
				},
				{
					ID:     "3",
					Type:   "A",
					Name:   "dns.hetzner.com",
					ZoneID: "2",
				},
			},
		}

		switch r.URL.Query().Get("zone_id") {
		case "1":
			resp.Records = []schema.Record{resp.Records[0]}
		case "3":
			resp.Records = []schema.Record{}
		}

		json.NewEncoder(w).Encode(resp) // nolint: errcheck
	})

	opts := RecordListOpts{}
	zones, _, err := env.Client.Record.List(env.Context, opts)
	as.NoError(err)
	as.EqInt(3, len(zones))

	opts.ZoneID = "1"
	zones, _, err = env.Client.Record.List(env.Context, opts)
	as.NoError(err)
	as.EqInt(1, len(zones))

	opts.ZoneID = "3"
	zones, _, err = env.Client.Record.List(env.Context, opts)
	as.NoError(err)
	as.EqInt(0, len(zones))
}

func TestRecordGetByID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	as := newAssert(t)

	env.Mux.HandleFunc(fmt.Sprintf("%s/1", pathRecords), func(w http.ResponseWriter, r *http.Request) {
		var resp schema.RecordResponse
		resp.Record = schema.Record{
			ID:   "1",
			Name: "hetzner.com",
		}

		json.NewEncoder(w).Encode(resp) // nolint: errcheck
	})

	id := "0"
	_, resp, err := env.Client.Record.GetByID(env.Context, id)
	if !as.EqInt(http.StatusNotFound, resp.StatusCode) {
		as.NoError(err)
	}

	id = "1"
	rec, _, err := env.Client.Record.GetByID(env.Context, id)
	if as.NoError(err) {
		as.EqStr(id, rec.ID)
		as.EqStr("hetzner.com", rec.Name)
	}
}

func TestRecordCreate(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	as := newAssert(t)

	env.Mux.HandleFunc(pathRecords, func(w http.ResponseWriter, r *http.Request) {
		var body schema.RecordCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		var resp schema.RecordResponse
		resp.Record = schema.Record{
			ID:   "1",
			Name: body.Name,
		}

		json.NewEncoder(w).Encode(resp) // nolint: errcheck
	})

	opts := RecordCreateOpts{
		Name:  "",
		Type:  RecordTypeA,
		Value: "10.0.0.0",
		Zone:  &Zone{ID: "1"},
	}

	_, _, err := env.Client.Record.Create(env.Context, opts)
	if as.Error(err) {
		as.EqStr("name required", err.Error())
	}

	opts.Name = "hetzner.com"
	rec, _, err := env.Client.Record.Create(env.Context, opts)
	if as.NoError(err) {
		as.EqStr("1", rec.ID)
		as.EqStr(opts.Name, rec.Name)
	}
}

func TestRecordUpdate(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	as := newAssert(t)

	env.Mux.HandleFunc(fmt.Sprintf("%s/1", pathRecords), func(w http.ResponseWriter, r *http.Request) {
		var body schema.RecordCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		var resp schema.RecordResponse
		resp.Record = schema.Record{
			ID:   "1",
			Name: body.Name,
		}

		json.NewEncoder(w).Encode(resp) // nolint: errcheck
	})

	rc := &Record{ID: "1"}
	opts := RecordUpdateOpts{
		Name:  "",
		Type:  RecordTypeA,
		Value: "10.0.0.0",
		Zone:  &Zone{ID: "1"},
	}

	_, _, err := env.Client.Record.Update(env.Context, rc, opts)
	if as.Error(err) {
		as.EqStr("name required", err.Error())
	}

	rc = &Record{ID: "0"}
	opts.Name = "dns.hetzner.com"
	_, resp, err := env.Client.Record.Update(env.Context, rc, opts)
	if !as.EqInt(http.StatusNotFound, resp.StatusCode) {
		as.NoError(err)
	}

	rc = &Record{ID: "1"}
	opts.Name = "dns.hetzner.com"
	recNew, _, err := env.Client.Record.Update(env.Context, rc, opts)
	if as.NoError(err) {
		as.EqStr(rc.ID, recNew.ID)
		as.EqStr(opts.Name, recNew.Name)
	}
}

func TestRecordDelete(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	as := newAssert(t)

	env.Mux.HandleFunc(fmt.Sprintf("%s/1", pathRecords), func(w http.ResponseWriter, r *http.Request) {})

	rec := &Record{ID: "0"}
	resp, err := env.Client.Record.Delete(env.Context, rec)
	if !as.EqInt(http.StatusNotFound, resp.StatusCode) {
		as.NoError(err)
	}

	rec.ID = "1"
	_, err = env.Client.Record.Delete(env.Context, rec)
	as.NoError(err)
}

func TestRecordBulkCreate(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	as := newAssert(t)

	env.Mux.HandleFunc(fmt.Sprintf("%s/bulk", pathRecords), func(w http.ResponseWriter, r *http.Request) {
		var resp schema.RecordBulkCreateResponse
		resp.Records = []schema.Record{
			{
				ID:     "1",
				Type:   "A",
				Name:   "hetzner.cloud",
				ZoneID: "1",
			},
			{
				ID:     "2",
				Type:   "A",
				Name:   "hetzner.com",
				ZoneID: "2",
			},
			{
				ID:     "3",
				Type:   "A",
				Name:   "dns.hetzner.com",
				ZoneID: "2",
			},
		}

		resp.ValidRecords = []schema.RecordBulkEntry{
			{
				Type:   "A",
				Name:   "hetzner.cloud",
				ZoneID: "1",
			},
		}

		resp.InvalidRecords = []schema.RecordBulkEntry{
			{
				Type:   "A",
				Name:   "hetzner.com",
				ZoneID: "2",
			},
			{
				Type:   "A",
				Name:   "dns.hetzner.com",
				ZoneID: "2",
			},
		}

		json.NewEncoder(w).Encode(resp) // nolint: errcheck
	})

	opts := []RecordCreateOpts{
		{
			Type:  RecordTypeA,
			Name:  "hetzner.cloud",
			Value: "",
			Zone:  &Zone{ID: "1"},
		},
		{
			Type:  RecordTypeA,
			Name:  "hetzner.com",
			Value: "10.0.0.2",
			Zone:  &Zone{ID: "2"},
		},
		{
			Type:  RecordTypeA,
			Name:  "dns.hetzner.com",
			Value: "10.0.0.2",
			Zone:  &Zone{ID: "2"},
		},
	}

	_, _, err := env.Client.Record.BulkCreate(env.Context, opts)
	if as.Error(err) {
		as.EqStr("value required", err.Error())
	}

	opts[0].Value = "10.0.0.1"

	br, _, err := env.Client.Record.BulkCreate(env.Context, opts)
	if as.NoError(err) {
		as.EqInt(3, len(br.Records))
		as.EqInt(1, len(br.ValidRecords))
		as.EqInt(2, len(br.InvalidRecords))
	}
}

func TestRecordBulkUpdate(t *testing.T) {
	//TODO
	env := newTestEnv()
	defer env.Teardown()

	as := newAssert(t)

	env.Mux.HandleFunc(fmt.Sprintf("%s/bulk", pathRecords), func(w http.ResponseWriter, r *http.Request) {
		var resp schema.RecordBulkUpdateResponse
		resp.Records = []schema.Record{
			{
				ID:     "1",
				Type:   "A",
				Name:   "hetzner.cloud",
				ZoneID: "1",
			},
			{
				ID:     "2",
				Type:   "A",
				Name:   "hetzner.com",
				ZoneID: "2",
			},
			{
				ID:     "3",
				Type:   "A",
				Name:   "dns.hetzner.com",
				ZoneID: "2",
			},
		}

		resp.FailedRecords = []schema.RecordBulkEntry{
			{
				Type:   "A",
				Name:   "hetzner.cloud",
				ZoneID: "1",
			},
		}

		json.NewEncoder(w).Encode(resp) // nolint: errcheck

	})

	opts := []RecordBulkUpdateOpts{
		{
			ID:    "",
			Type:  RecordTypeA,
			Name:  "hetzner.cloud",
			Value: "10.0.0.1",
			Zone:  &Zone{ID: "1"},
		},
		{
			ID:    "2",
			Type:  RecordTypeA,
			Name:  "hetzner.com",
			Value: "10.0.0.2",
			Zone:  &Zone{ID: "2"},
		},
		{
			ID:    "3",
			Type:  RecordTypeA,
			Name:  "dns.hetzner.com",
			Value: "10.0.0.2",
			Zone:  &Zone{ID: "2"},
		},
	}

	_, _, err := env.Client.Record.BulkUpdate(env.Context, opts)
	if as.Error(err) {
		as.EqStr("id required", err.Error())
	}

	opts[0].ID = "1"
	br, _, err := env.Client.Record.BulkUpdate(env.Context, opts)
	if as.NoError(err) {
		as.EqInt(3, len(br.Records))
		as.EqInt(1, len(br.FailedRecords))
	}
}
