package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type results struct {
	mu *sync.Mutex
	s  []string
}

var res = results{
	mu: &sync.Mutex{},
	s:  []string{},
}

func main() {

	/*
			client := NewClient(5, 100, time.Second*5, time.Second)
			time.Sleep(time.Second * 2)

			for i := 0; i < 10; i++ {
				go client.LimitedRequests(i)

			}

			time.Sleep(time.Second * 5)


		t := time.Second / 2

		limit := NewLimiter(1, 10, time.Second*2)
		go lineup(limit, nil, 1)
		time.Sleep(time.Second / 7)
		go lineup(limit, nil, 2)
		time.Sleep(time.Second / 7)
		go lineup(limit, &t, 3)
		time.Sleep(time.Second / 7)
		go lineup(limit, nil, 4)
		time.Sleep(time.Second / 7)
		go lineup(limit, nil, 5)
		time.Sleep(time.Second / 7)
		go lineup(limit, nil, 6)
		time.Sleep(time.Second / 7)
		go lineup(limit, nil, 7)
		time.Sleep(time.Second / 7)
		go lineup(limit, nil, 8)
		time.Sleep(time.Second / 7)
		go lineup(limit, nil, 9)
		time.Sleep(time.Second / 7)
		go lineup(limit, nil, 10)
		time.Sleep(time.Second / 7)
		go lineup(limit, nil, 11)
		time.Sleep(time.Second / 7)
		go lineup(limit, nil, 12)
		time.Sleep(time.Second / 7)
		go lineup(limit, nil, 13)
		time.Sleep(time.Second * 30)

		for _, str := range res.s {
			fmt.Print(str)
		}
	*/

	l := NewUnbufferedLimiter(5, time.Second)
	n := time.Now()
	for i := 0; i < 50; i++ {
		err := l.Wait(nil)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			time.Sleep(time.Second / 10)
			continue
		}
		fmt.Printf("for %d, time passed: %v\n", i, time.Since(n))
		time.Sleep(time.Second / 10)
	}

}

func lineup(lim *Limiter, t *time.Duration, thread int) {
	res.mu.Lock()
	res.s = append(res.s, fmt.Sprintf("started thread %d\n", thread))
	res.mu.Unlock()
	err := lim.Wait(t)
	if err != nil {
		res.mu.Lock()
		res.s = append(res.s, fmt.Sprintf("error from thread %d: %v\n", thread, err))
		res.mu.Unlock()
		return
	}
	res.mu.Lock()
	res.s = append(res.s, fmt.Sprintf("processed request from thread %d\n", thread))
	res.mu.Unlock()
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
