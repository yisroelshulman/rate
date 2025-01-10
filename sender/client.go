package main

import (
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	ClientLimiter *Limiter
	HttpClient    http.Client
}

func NewClient(requestsAllowed int, timeout, interval time.Duration) *Client {
	client := &Client{
		ClientLimiter: NewLimiter(requestsAllowed, interval),
		HttpClient: http.Client{
			Timeout: timeout,
		},
	}
	return client
}

func (c *Client) DoRequest(req *http.Request) {
	err := c.ClientLimiter.Wait(nil)
	if err != nil {
		fmt.Printf("failed")
		return
	}
	_, err = c.HttpClient.Do(req)
	if err != nil {
		return
	}
}
