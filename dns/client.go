package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/jobstoit/hetzner-dns-go/dns/schema"
)

var validTokenReg = regexp.MustCompile("[a-zA-Z0-9]{32}")

// Client is the client for the Hetzner DNS API.
type Client struct {
	httpClient         *http.Client
	token              string
	endpoint           string
	debugWriter        io.Writer
	tokenValid         bool
	applicationName    string
	applicationVersion string
	userAgent          string

	Zone          *ZoneClient
	Record        *RecordClient
	PrimaryServer *PrimaryServerClient
}

// ClientOption is used to configure a client.
type ClientOption func(*Client)

// WithToken configures a client to use the specified token.
func WithToken(token string) ClientOption {
	return func(client *Client) {
		client.token = token
		client.tokenValid = validTokenReg.MatchString(token)
	}
}

// WithEndpoint configures the client to use a different endpoint.
func WithEndpoint(url string) ClientOption {
	return func(client *Client) {
		client.endpoint = strings.TrimRight(url, "/")
	}
}

// WithApplication configures a Client with the given application name and
// application version. The version may be blank. Programs are encouraged
// to at least set an application name.
func WithApplication(name, version string) ClientOption {
	return func(client *Client) {
		client.applicationName = name
		client.applicationVersion = version
	}
}

// WithHTTPClient configures a Client to perform HTTP requests with httpClient.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

// WithDebugWriter configures the client to use the given debug writer.
func WithDebugWriter(wr io.Writer) ClientOption {
	return func(client *Client) {
		client.debugWriter = wr
	}
}

// NewClient creates a new client.
func NewClient(options ...ClientOption) *Client {
	client := &Client{
		endpoint:   Endpoint,
		tokenValid: true,
		httpClient: http.DefaultClient,
	}

	for _, option := range options {
		option(client)
	}

	client.buildUserAgent()

	client.Zone = &ZoneClient{client}
	client.Record = &RecordClient{client}
	client.PrimaryServer = &PrimaryServerClient{client}

	return client
}

// NewRequest creates an HTTP request against the API. The returned request
// is assigned with ctx and has all necessary headers set (auth, user agent, etc.).
func (c *Client) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	url := c.endpoint + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)

	if !c.tokenValid {
		return nil, errors.New("authorization token contains invalid characters")
	} else if c.token != "" {
		req.Header.Set("Auth-API-Token", c.token)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req = req.WithContext(ctx)
	return req, nil
}

func (c *Client) buildUserAgent() {
	switch {
	case c.applicationName != "" && c.applicationVersion != "":
		c.userAgent = c.applicationName + "/" + c.applicationVersion + " " + UserAgent
	case c.applicationName != "" && c.applicationVersion == "":
		c.userAgent = c.applicationName + " " + UserAgent
	default:
		c.userAgent = UserAgent
	}
}

// Do performs an HTTP request against the API.
func (c *Client) Do(r *http.Request, v interface{}) (*Response, error) {
	var body []byte
	var err error
	if r.ContentLength > 0 {
		body, err = io.ReadAll(r.Body)
		if err != nil {
			r.Body.Close()
			return nil, err
		}
		r.Body.Close()
	}

	if r.ContentLength > 0 {
		r.Body = io.NopCloser(bytes.NewReader(body))
	}

	if c.debugWriter != nil {
		dumpReq, err := dumpRequest(r)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(c.debugWriter, "--- Request:\n%s\n\n", dumpReq)
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	response := &Response{Response: resp}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return response, err
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(body))

	if c.debugWriter != nil {
		dumpResp, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(c.debugWriter, "--- Response:\n%s\n\n", dumpResp)
	}

	if err = response.readMeta(body); err != nil {
		return response, fmt.Errorf("hetzner-dns: error reading response meta data: %s", err)
	}

	if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
		err = fmt.Errorf("hetzner-dns: server responded with status code %d", resp.StatusCode)
		return response, err
	}
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, bytes.NewReader(body))
		} else {
			err = json.Unmarshal(body, v)
		}
	}

	return response, err
}

// ListOpts specifies options for listing resources
type ListOpts struct {
	Page    int
	PerPage int
}

func (l ListOpts) values() url.Values {
	vals := url.Values{}
	if l.Page > 0 {
		vals.Add("page", strconv.Itoa(l.Page))
	}
	if l.PerPage > 0 {
		vals.Add("per_page", strconv.Itoa(l.PerPage))
	}
	return vals
}

// Response represents a response from the API. It embeds http.Response.
type Response struct {
	*http.Response
	Meta Meta
}

func (r *Response) readMeta(body []byte) error {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		var s schema.MetaResponse
		if err := json.Unmarshal(body, &s); err != nil {
			return err
		}
		if s.Meta.Pagination != nil {
			p := PaginationFromSchema(*s.Meta.Pagination)
			r.Meta.Pagination = &p
		}
	}

	return nil
}

// Meta represents meta inforation included in the API response.
type Meta struct {
	Pagination *Pagination
}

// Pagination represents pagination meta information.
type Pagination struct {
	LastPage     int
	Page         int
	PerPage      int
	TotalEntries int
}

func dumpRequest(r *http.Request) ([]byte, error) {
	// Duplicate the request, so we can redact the auth header
	rDuplicate := r.Clone(context.Background())
	rDuplicate.Header.Set("Authorization", "REDACTED")

	// To get the request body we need to read it before the request was actually sent.
	// See https://github.com/golang/go/issues/29792
	dumpReq, err := httputil.DumpRequestOut(rDuplicate, true)
	if err != nil {
		return nil, err
	}

	// Set original request body to the duplicate created by DumpRequestOut. The request body is not duplicated
	// by .Clone() and instead just referenced, so it would be completely read otherwise.
	r.Body = rDuplicate.Body

	return dumpReq, nil
}
