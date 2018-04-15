package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
)

const v1 = "/api/v1"

// Client wraps http client.
type Client struct {
	httpClient *http.Client
	endpoint   *url.URL
	apiKey     string
	secretKey  string
	async      bool
	timeout    time.Duration
}

func newClient(httpClient *http.Client, endpoint, apiKey, secretKey string, async bool) (*Client, error) {
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	return &Client{
		httpClient: httpClient,
		endpoint:   endpointURL,
		apiKey:     apiKey,
		secretKey:  secretKey,
		async:      async,
		timeout:    60 * time.Second,
	}, nil
}

// NewClient returns new Client.
func NewClient(httpClient *http.Client, endpoint, apiKey, secretKey string) (*Client, error) {
	return newClient(httpClient, endpoint, apiKey, secretKey, false)
}

// NewAsyncClient returns new Client.
func NewAsyncClient(httpClient *http.Client, endpoint, apiKey, secretKey string) (*Client, error) {
	return newClient(httpClient, endpoint, apiKey, secretKey, true)
}

// AsyncTimeout set a timeout second for async job.
func (c *Client) AsyncTimeout(timeout time.Duration) {
	c.timeout = timeout
}

func (c *Client) newRequest(ctx context.Context, method, spath string, query []byte) (*http.Request, error) {
	endpoint := *c.endpoint
	endpoint.Path = path.Join(endpoint.Path, spath)

	req, err := http.NewRequest(
		method,
		endpoint.String(),
		bytes.NewBuffer(query),
	)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	return req, nil
}

// WaitAsyncJob is a helper function that waits a async job to done.
func (c *Client) WaitAsyncJob(jobID string) error {
	var timeout <-chan time.Time
	timeout = time.After(c.timeout)

	var lastErr error
	for {
		select {
		case <-timeout:
			return fmt.Errorf("JobID \"%s\" timeouted. last error: %s", jobID, lastErr)
		default:
			ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
			defer cancel()

			r, err := c.GetJob(ctx, jobID)
			if r.JobStatus == "End" {
				return nil
			}
			if lastErr = err; err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
		}
	}
}
