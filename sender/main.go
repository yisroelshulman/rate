package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	client := NewClient(50, time.Second*5, time.Second)
	go client.LimitedRequests()
	go client.LimitedRequests()
	go client.LimitedRequests()
	go client.LimitedRequests()
	go client.LimitedRequests()
	go client.LimitedRequests()
	go client.LimitedRequests()
	go client.LimitedRequests()

	time.Sleep(time.Second * 10)
}

func UnlimitedRequests() {
	client := http.Client{}
	url := "http://localhost:9578/healthz"
	for {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("failed with: %v\n", err)
			return
		}
		_, err = client.Do(req)
		if err != nil {
			fmt.Printf("req errd with: %v\n", err)
		}
	}
}

func (c *Client) LimitedRequests() {
	url := "http://localhost:9578/healthz"
	for {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("failed with: %v\n", err)
			return
		}
		c.Request(req)
	}
}
