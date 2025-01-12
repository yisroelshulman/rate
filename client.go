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

func NewClient(requestsAllowed, capacity int, timeout, interval time.Duration) *Client {
	client := &Client{
		ClientLimiter: NewLimiter(requestsAllowed, capacity, interval),
		HttpClient: http.Client{
			Timeout: timeout,
		},
	}
	return client
}

func (c *Client) DoRequest(req *http.Request, threadNum, requestNum int) {
	//fmt.Printf("Request from %d, requst num %d >>>\n", threadNum, requestNum)
	err := c.ClientLimiter.Wait(nil)
	if err != nil {
		fmt.Printf("	== timeout from %d, requst num %d\n", threadNum, requestNum)
		//fmt.Printf("failed")
		return
	}
	//fmt.Printf("<<< Processed from %d, requst num %d\n", threadNum, requestNum)
	_, err = c.HttpClient.Do(req)
	if err != nil {
		return
	}
}
