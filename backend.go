package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// AddParam is the request parameter used by the API.
type AddParam struct {
	IPAddress string `json:"ipaddress"`
	Port      int    `json:"port"`
}

// ListResponse is the response data returned from the API.
type ListResponse []struct {
	ID   string `json:"id"`
	Name int    `json:"name"`
}

// Response is the response data returned from the API.
type Response struct {
	JobID  string `json:"job_id"`
	Status int    `json:"job_status"`
}

// List calls the API that fetches list of something.
func (c *Client) List(ctx context.Context) (*ListResponse, error) {
	path := "/path"

	httpReq, err := c.newRequest(ctx, "GET", path, nil)
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

	var res = new(ListResponse)
	if err := decodeBody(httpRes, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Add calls the API that adds something.
func (c *Client) Add(ctx context.Context, id string, p *AddParam) (*Response, error) {
	path := fmt.Sprintf("/path/%s", id)

	query, _ := json.Marshal(p)

	httpReq, err := c.newRequest(ctx, "POST", path, query)
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

	var res = new(Response)
	if err := decodeBody(httpRes, res); err != nil {
		return nil, err
	}

	if c.async {
		if err := c.WaitAsyncJob(res.JobID); err != nil {
			return nil, err
		}
	}

	return res, nil
}

// Delete calls the API that deletes something.
func (c *Client) Delete(ctx context.Context, id string) (*Response, error) {
	path := fmt.Sprintf("/path/%s", id)

	httpReq, err := c.newRequest(ctx, "DELETE", path, nil)
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

	var res = new(Response)
	if err := decodeBody(httpRes, res); err != nil {
		return nil, err
	}

	if c.async {
		if err := c.WaitAsyncJob(res.JobID); err != nil {
			return nil, err
		}
	}

	return res, nil
}
