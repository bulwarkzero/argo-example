package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Response struct {
	Status int
	Code   string
	Data   interface{}
}

const (
	httpRequestsCount    = "requests_total"
	httpRequestsDuration = "request_duration_seconds"
)

var buckets = []float64{
	0.0005,
	0.001, // 1ms
	0.002,
	0.005,
	0.01, // 10ms
	0.02,
	0.05,
	0.1, // 100 ms
	0.2,
	0.5,
	1.0, // 1s
	2.0,
	5.0,
	10.0, // 10s
	15.0,
	20.0,
	30.0,
}

var httpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "troll",
	Name:      httpRequestsCount,
	Help:      "HTTP Request count",
}, []string{"status", "method", "handler"})

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "troll",
	Name:      httpRequestsDuration,
	Help:      "Duration of each request",
	Buckets:   buckets,
}, []string{"method", "handler"})

func handleHello(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(httpDuration.WithLabelValues(r.Method, r.URL.Path))

	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)

	response := &Response{
		Status: http.StatusOK,
		Code:   "Ok",
		Data: map[string]interface{}{
			"msg": "Hello from test server gada",
		},
	}

	rawResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))

		timer.ObserveDuration()
		httpRequests.WithLabelValues("500", r.Method, r.URL.Path).Inc()

		return
	}

	timer.ObserveDuration()
	httpRequests.WithLabelValues("200", r.Method, r.URL.Path).Inc()

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(rawResponse); err != nil {
		log.Println("error on writing response", err)
	}
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/hello", handleHello)
	http.HandleFunc("/healthz", handleHealthz)

	addr := os.Getenv("APP_SERVE_ADDR")
	if addr == "" {
		addr = ":5040"
	}

	log.Fatal(http.ListenAndServe(addr, nil))
}
