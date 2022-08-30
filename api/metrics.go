package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type MetricsServer struct {
	Logger *zap.SugaredLogger
	Server *http.Server
}

func NewMetricsServer(address string, logger *zap.SugaredLogger) *MetricsServer {
	return &MetricsServer{
		Server: &http.Server{
			Addr:         address,
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      loggingMiddleware(promhttp.Handler()),
		},
		Logger: logger,
	}
}
