package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	client := NewClient(100, time.Second*5, time.Second)
	for i := 0; i < 100000; i++ {
		go client.LimitedRequests()
	}

	time.Sleep(time.Second * 100)

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

		c.DoRequest(req)
	}
}
