package main

import (
	_ "expvar"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
)

type DebugServer struct {
	Logger *zap.SugaredLogger
	Server *http.Server
}

func NewDebugServer(address string, logger *zap.SugaredLogger) *DebugServer {
	return &DebugServer{
		Server: &http.Server{
			Addr:    address,
			Handler: loggingMiddleware(http.DefaultServeMux),
		},
		Logger: logger,
	}
}
