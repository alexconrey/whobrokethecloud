package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
)

type MetricsServer struct {
	Logger *zap.SugaredLogger
	Server *http.Server
}

func NewMetricsServer(address string, logger *zap.SugaredLogger) *MetricsServer {
	return &MetricsServer{
		Server: &http.Server{
			Addr:    address,
			Handler: loggingMiddleware(promhttp.Handler()),
		},
		Logger: logger,
	}
}
