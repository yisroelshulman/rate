package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type config struct {
	requestHits atomic.Int32
}

func main() {

	cfg := config{
		requestHits: atomic.Int32{},
	}

	port := "9578"
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", cfg.HandlerReadiness)
	mux.HandleFunc("GET /metrics", cfg.HandlerMetrics)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("starting server on port %s\n", port)
	server.ListenAndServe()

}

func (cfg *config) HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	cfg.requestHits.Add(1)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
	r.Body.Close()
}

func (cfg *config) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome to the Test Metrics</h1>
				<p>Total requests made: %d</p>
			</body>
		</htlml>
	`, cfg.requestHits.Load())))
}
