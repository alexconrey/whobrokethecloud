package main

import (
	_ "expvar"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type DebugServer struct {
	Logger *zap.SugaredLogger
	Server *http.Server
}

func NewDebugServer(address string, logger *zap.SugaredLogger) *DebugServer {
	return &DebugServer{
		Server: &http.Server{
			Addr:         address,
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      loggingMiddleware(http.DefaultServeMux),
		},
		Logger: logger,
	}
}
