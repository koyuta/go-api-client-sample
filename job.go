package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// GetJobResponse is the response data returned from calling the job API.
type GetJobResponse struct {
	ID     string `json:"job_id"`
	Status string `json:"job_status"`
}

// GetJob gets job data by jobID.
func (c *Client) GetJob(ctx context.Context, jobID string) (*GetJobResponse, error) {
	path := fmt.Sprintf("/jobs/%s", jobID)

	httpReq, err := c.newRequest(ctx, "GET", path, nil, time.Time{})
	if err != nil {
		return nil, err
	}

	httpRes, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode != http.StatusOK {
	}

	b, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	var res = new(GetJobResponse)
	if err := decodeBody(httpRes, res); err != nil {
		return nil, err
	}

	return res, nil
}
