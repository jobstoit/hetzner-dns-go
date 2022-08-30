package dns

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jobstoit/hetzner-dns-go/dns/schema"
)

func TestPrimaryServerList(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	as := newAssert(t)

	env.Mux.HandleFunc(pathPrimaryServers, func(w http.ResponseWriter, r *http.Request) {
		var respBody schema.PrimaryServerListResponse
		respBody.PrimaryServers = []schema.PrimaryServer{
			{
				ID: "1",
			},
			{
				ID: "2",
			},
			{
				ID: "3",
			},
		}

		if r.URL.Query().Get("zone_id") == "2" {
			respBody.PrimaryServers = []schema.PrimaryServer{
				respBody.PrimaryServers[0],
			}
		}

		json.NewEncoder(w).Encode(respBody) // nolint: errcheck
	})

	opts := PrimaryServerListOpts{}
	svrs, _, err := env.Client.PrimaryServer.List(env.Context, opts)
	if as.NoError(err) {
		as.EqInt(3, len(svrs))
	}

	opts.ZoneID = "2"
	svrs, _, err = env.Client.PrimaryServer.List(env.Context, opts)
	if as.NoError(err) {
		as.EqInt(1, len(svrs))
	}
}

func TestPrimaryServerGetByID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	as := newAssert(t)

	env.Mux.HandleFunc(fmt.Sprintf("%s/1", pathPrimaryServers), func(w http.ResponseWriter, r *http.Request) {
		var respBody schema.PrimaryServerResponse
		respBody.PrimaryServer = schema.PrimaryServer{
			ID: "1",
		}

		json.NewEncoder(w).Encode(respBody) // nolint: errcheck
	})

	svr, _, err := env.Client.PrimaryServer.GetByID(env.Context, "1")
	if as.NoError(err) {
		as.EqStr("1", svr.ID)
	}
}

func TestPrimaryServerCreate(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	as := newAssert(t)

	env.Mux.HandleFunc(pathPrimaryServers, func(w http.ResponseWriter, r *http.Request) {
		var body schema.PrimaryServerUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var respBody schema.PrimaryServerResponse
		respBody.PrimaryServer = schema.PrimaryServer{
			ID:       "1",
			Port:     body.Port,
			Created:  time.Now(),
			Modified: time.Now(),
			ZoneID:   body.ZoneID,
			Address:  body.Address,
		}

		json.NewEncoder(w).Encode(respBody) // nolint: errcheck
	})

	opts := PrimaryServerCreateOpts{
		Address: "",
		Port:    80,
		ZoneID:  "2",
	}

	_, _, err := env.Client.PrimaryServer.Create(env.Context, opts)
	if as.Error(err) {
		as.EqStr("address required", err.Error())
	}

	opts.Address = "dns.hetzner.com"
	svr, _, err := env.Client.PrimaryServer.Create(env.Context, opts)
	if as.NoError(err) {
		as.EqStr("1", svr.ID)
		as.EqStr(opts.Address, svr.Address)
	}
}

func TestPrimaryServerUpdate(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	as := newAssert(t)

	env.Mux.HandleFunc(fmt.Sprintf("%s/1", pathPrimaryServers), func(w http.ResponseWriter, r *http.Request) {
		var body schema.PrimaryServerUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var respBody schema.PrimaryServerResponse
		respBody.PrimaryServer = schema.PrimaryServer{
			ID:       "1",
			Port:     body.Port,
			Created:  time.Now(),
			Modified: time.Now(),
			ZoneID:   body.ZoneID,
			Address:  body.Address,
		}

		json.NewEncoder(w).Encode(respBody) // nolint: errcheck
	})

	ps := &PrimaryServer{ID: "0"}
	opts := PrimaryServerUpdateOpts{
		Address: "dns.hetzner.com",
		Port:    80,
		ZoneID:  "2",
	}

	_, resp, err := env.Client.PrimaryServer.Update(env.Context, ps, opts)
	if !as.EqInt(http.StatusNotFound, resp.StatusCode) {
		as.NoError(err)
	}

	ps.ID = "1"
	opts.Address = ""
	_, _, err = env.Client.PrimaryServer.Update(env.Context, ps, opts)
	if as.Error(err) {
		as.EqStr("address required", err.Error())
	}

	opts.Address = "dns.hetzner.com"
	svr, _, err := env.Client.PrimaryServer.Update(env.Context, ps, opts)
	if as.NoError(err) {
		as.EqStr("1", svr.ID)
		as.EqStr(opts.Address, svr.Address)
	}
}

func TestPrimaryServerDelete(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	as := newAssert(t)

	env.Mux.HandleFunc(fmt.Sprintf("%s/1", pathPrimaryServers), func(w http.ResponseWriter, r *http.Request) {})

	ps := &PrimaryServer{ID: "0"}
	resp, err := env.Client.PrimaryServer.Delete(env.Context, ps)
	if !as.EqInt(http.StatusNotFound, resp.StatusCode) {
		as.NoError(err)
	}

	ps.ID = "1"
	_, err = env.Client.PrimaryServer.Delete(env.Context, ps)
	as.NoError(err)
}
