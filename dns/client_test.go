package dns

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testEnv struct {
	Server  *httptest.Server
	Mux     *http.ServeMux
	Client  *Client
	Context context.Context
}

func (env *testEnv) Teardown() {
	env.Server.Close()
	env.Server = nil
	env.Mux = nil
	env.Client = nil
	env.Context.Done()
	env.Context = nil
}

func newTestEnv() testEnv {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := NewClient(
		WithEndpoint(server.URL),
		WithToken("32CharactersTokenxxxxxxxXxxxxxxx"),
		WithApplication("testing", Version),
	)
	return testEnv{
		Server:  server,
		Mux:     mux,
		Client:  client,
		Context: context.Background(),
	}
}

func TestClientEndpointTrailingSlashesRemoved(t *testing.T) {
	client := NewClient(WithEndpoint("http://api/v1.0/////"))
	if strings.HasSuffix(client.endpoint, "/") {
		t.Fatalf("endpoint has trailing slashes: %q", client.endpoint)
	}
}

func TestClientInvalidToken(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Client = NewClient(
		WithEndpoint(env.Server.URL),
		WithToken("32Charact3rsInvalidT@k3n!xxxxxxx"),
	)

	ctx := context.Background()
	_, err := env.Client.NewRequest(ctx, "GET", "/", nil)

	if nil == err {
		t.Error("Failed to trigger expected error")
	} else if err.Error() != "authorization token contains invalid characters" {
		t.Fatalf("Invalid encoded authorization token triggered unexpected error message: %s", err)
	}
}

func TestClientDo(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	callCount := 1
	env.Mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		switch callCount {
		case 1:
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintln(w, "{}")
		case 2:
			fmt.Fprintln(w, "{}")
		default:
			t.Errorf("unexpected number of calls to the test server: %v", callCount)
		}
	})

	ctx := context.Background()
	request, _ := env.Client.NewRequest(ctx, http.MethodGet, "/test", nil)
	_, err := env.Client.Do(request, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 2 {
		t.Fatalf("unexpected callCount: %v", callCount)
	}
}
