package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"os"
	"strconv"
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"handler", "code"})

	requestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"handler", "code"})
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

// prometheusMiddleware implements mux.MiddlewareFunc.
func prometheusMiddleware(next http.Handler) http.Handler {
	prometheus.Register(httpDuration)
	prometheus.Register(requestCounter)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		wr := NewLoggingResponseWriter(w)
		statusCode := strconv.Itoa(wr.statusCode)

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path, statusCode))
		next.ServeHTTP(wr, r)
		fmt.Println(wr.statusCode)
		timer.ObserveDuration()
		requestCounter.WithLabelValues(path, statusCode).Inc()
	})
}

func (o *Outages) Handler(w http.ResponseWriter, r *http.Request) {
	var outages = make(map[string][]Outage)
	for _, outage := range o.Outages {
		provider := outage.Provider
		providerMap, ok := outages[provider]
		if !ok {
			outages[provider] = make([]Outage, 0)
		}
		outages[provider] = append(providerMap, outage)
	}

	data, err := json.Marshal(outages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Write(data)
}
