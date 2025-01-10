package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type Client struct {
	CycleRequestCount atomic.Int32
	RequestsPerCycle  int
	HttpClient        http.Client
}

func NewClient(requestsAllowed int, timeout, interval time.Duration) *Client {
	client := &Client{
		CycleRequestCount: atomic.Int32{},
		RequestsPerCycle:  requestsAllowed,
		HttpClient: http.Client{
			Timeout: timeout,
		},
	}

	go client.resetCountLoop(interval)

	return client
}

func (c *Client) resetCountLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reset()
	}
}

func (c *Client) reset() {
	a := c.CycleRequestCount.Load()
	fmt.Printf("reset from: %d\n", a)
	c.CycleRequestCount.Store(0)
}

func (c *Client) Request(req *http.Request) {
	waiting := false
	for c.CycleRequestCount.Load() >= int32(c.RequestsPerCycle) {
		if !waiting {
			fmt.Println("waiting...")
		}
		waiting = true
		continue
	}
	c.CycleRequestCount.Add(1)
	_, err := c.HttpClient.Do(req)
	if err != nil {
		return
	}

}
