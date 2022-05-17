package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/svanharmelen/jsonapi"
	"golang.org/x/time/rate"
)

var (
	// ErrBadRequest is returned when a receiving a 400.
	ErrBadRequest = errors.New("bad request")
	// ErrResourceNotFound is returned when a receiving a 404.
	ErrResourceNotFound = errors.New("resource not found")
	// ErrGatewayTimeout is returned when a receiving a 504.
	ErrGatewayTimeout = errors.New("gateway timeout")
)

func DefaultConfig() *Config {
	const (
		userAgent       = "petstore-sdk-v1"
		DefaultAddress  = "http://localhost:8080"
		DefaultBasePath = "/api/v1"
	)

	config := &Config{
		BasePath:   DefaultBasePath,
		Headers:    http.Header{},
		HTTPClient: cleanhttp.DefaultPooledClient(),
		Logger:     log.Default(),
	}

	config.Address = os.Getenv("PETSTORE_ADDRESS")
	if config.Address == "" {
		config.Address = DefaultAddress
	}

	config.Headers.Set("User-Agent", userAgent)
	config.Headers.Set("Content-Type", "application/json")

	return config
}

type Config struct {
	Address    string
	BasePath   string
	Headers    http.Header
	HTTPClient *http.Client
	Logger     *log.Logger
}

func NewClient(cfg *Config) (*Client, error) {
	config := DefaultConfig()
	if cfg != nil {
		if cfg.Address != "" {
			config.Address = cfg.Address
		}
		if cfg.BasePath != "" {
			config.BasePath = cfg.BasePath
		}
		for k, v := range cfg.Headers {
			config.Headers[k] = v
		}
		if cfg.HTTPClient != nil {
			config.HTTPClient = cfg.HTTPClient
		}
		if cfg.Logger != nil {
			config.Logger = cfg.Logger
		}
	}

	baseURL, err := url.Parse(config.Address + config.BasePath)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %v", err)
	}

	client := &Client{
		baseURL: baseURL,
		headers: config.Headers,
		http:    config.HTTPClient,
		log:     config.Logger,
	}

	// Create resources...
	client.Pets = newPets(client)

	return client, nil
}

type Client struct {
	baseURL *url.URL
	headers http.Header
	http    *http.Client
	log     *log.Logger
	limiter *rate.Limiter

	Pets Pets
}

func (c *Client) debug(format string, v ...interface{}) {
	c.log.Printf("[DEBUG] go-petstore "+format, v...)
}

func (c *Client) newRequest(method, path string, v interface{}) (*http.Request, error) {
	u, err := url.Parse(c.baseURL.String() + "/" + path)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	switch method {
	case "GET":
		if v != nil {
			q, err := query.Values(v)
			if err != nil {
				return nil, err
			}
			u.RawQuery = q.Encode()
		}
	case "PATCH", "POST":
		if v != nil {
			d, _ := json.Marshal(v)
			c.debug("body: " + string(d))
			body = bytes.NewReader(d)
		}
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	for k, v := range c.headers {
		req.Header[k] = v
	}

	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) error {
	req = req.WithContext(ctx)
	c.debug("request: %v", req)

	resp, err := c.http.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return err
		}
	}
	defer resp.Body.Close()

	c.debug("response: %+v", resp)

	if err := c.checkResponseCode(resp); err != nil {
		if err != ErrGatewayTimeout {
			return err
		}
		// If the request timed out, retry once.
		resp, err = c.http.Do(req)
		if err := c.checkResponseCode(resp); err != nil {
			return err
		}
	}

	if v == nil {
		return nil
	}
	buf := bytes.Buffer{}
	buf.ReadFrom(resp.Body)
	return json.Unmarshal(buf.Bytes(), v)
}

// checkResponseCode can be used to check the status code of an HTTP request.
func (c *Client) checkResponseCode(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	}
	switch r.StatusCode {
	case 400:
		return ErrBadRequest
	case 404:
		return ErrResourceNotFound
	case 504:
		return ErrGatewayTimeout
	}

	// Decode the error payload.
	errPayload := &jsonapi.ErrorsPayload{}
	err := json.NewDecoder(r.Body).Decode(errPayload)
	if err != nil || len(errPayload.Errors) == 0 {
		c.debug("resp status: %+v", r.Status)
		return fmt.Errorf(r.Status)
	}

	// Parse and format the errors.
	var errs []string
	for _, e := range errPayload.Errors {
		if e.Detail == "" {
			errs = append(errs, e.Title)
		} else {
			errs = append(errs, fmt.Sprintf("%s\n\n%s", e.Title, e.Detail))
		}
	}

	return fmt.Errorf(strings.Join(errs, "\n"))
}
