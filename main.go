package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	client := NewClient(5, 100, time.Second*5, time.Second)
	time.Sleep(time.Second * 2)

	for i := 0; i < 10; i++ {
		go client.LimitedRequests(i)

	}

	time.Sleep(time.Second * 5)

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

func (c *Client) LimitedRequests(threadNum int) {
	url := "http://localhost:9578/healthz"
	for i := 0; i < 10; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("failed with: %v\n", err)
			return
		}

		c.DoRequest(req, threadNum, i)
	}
}
