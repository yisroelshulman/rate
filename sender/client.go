package main

import (
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	Limiter    *limiter
	HttpClient http.Client
}

func NewClient(requestsAllowed int, timeout, interval time.Duration) *Client {
	client := &Client{
		Limiter: NewLimiter(requestsAllowed, interval),
		HttpClient: http.Client{
			Timeout: timeout,
		},
	}

	return client

}

func (c *Client) RequestTest(req *http.Request) {
	err := c.Limiter.Wait(nil)
	if err != nil {
		fmt.Printf("failed")
		return
	}
	_, err = c.HttpClient.Do(req)
	if err != nil {
		return
	}
	return
}
