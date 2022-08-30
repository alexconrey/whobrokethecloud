package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	var (
		wait        time.Duration
		httpPort    int
		debugPort   int
		metricsPort int
		pollDelay   time.Duration
		isDebug     bool

		feedPollDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "outage_feed_poll_duration",
			Help: "Histogram for the runtime of a feed poll.",
		}, []string{"url"})
	)

	//
	// FLAGS
	//
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.IntVar(&httpPort, "http-port", 8080, "The port on which the HTTP server will run")
	flag.IntVar(&debugPort, "debug-port", 6060, "The port on which to run debugging services")
	flag.IntVar(&metricsPort, "metrics-port", 9100, "The port on which the metrics server will run")
	flag.DurationVar(&pollDelay, "poll-delay", time.Second*180, "The delay interval (in seconds) for waiting between polling attempts per provider")
	flag.BoolVar(&isDebug, "debug", false, "Enable debug features (pprof, etc)")
	flag.Parse()

	//
	// LOGGING
	//
	loggerConfig := zap.NewProductionConfig()

	if isDebug {
		loggerConfig.Level.SetLevel(zap.DebugLevel)
	}

	logger, _ := loggerConfig.Build()
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Info("Starting whobrokethe.cloud API service")

	// prometheus.Register(requestCounter)
	prometheus.Register(feedPollDuration)

	outages := NewOutages(sugar)

	// FEEDS
	aznFeeds := AmazonFeeds{
		Logger: sugar,
		Products: []string{
			"ec2",
			"ecs",
			"eks",
			"elb",
			"elasticache",
			"rds",
		},
		Regions: []string{
			"us-east-1",
			"us-east-2",
			"us-west-1",
			"us-west-2",
		},
		Chan:                  outages.Chan,
		PollDurationHistogram: feedPollDuration,
	}

	googleFeed := GoogleFeed{
		URL:                   "https://status.cloud.google.com/incidents.json",
		Logger:                sugar,
		Chan:                  outages.Chan,
		PollDurationHistogram: feedPollDuration,
	}

	// Polling goroutines
	go aznFeeds.Poll(pollDelay)
	go googleFeed.Poll(pollDelay)

	// Watch for outage events in the channel
	go outages.HandleOutages()

	// Poll metrics around outages
	go outages.PollMetrics()

	// Cleanup events in memory older than x amount of days
	go outages.StartWatchdog(90)

	// Start HTTP server
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Use(prometheusMiddleware)
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/outages", outages.Handler)
	r.HandleFunc("/healthz", HealthHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", httpPort),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Run HTTP server as non-blocking, goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	// Run Debug HTTP server as non-blocking, goroutine
	if isDebug {
		debugServer := NewDebugServer(fmt.Sprintf("%s:%d", "0.0.0.0", debugPort), sugar)
		sugar.Debug("Starting debug http server")
		go func() {
			if err := debugServer.Server.ListenAndServe(); err != nil {
				sugar.Error(err.Error())
			}
		}()
	}

	// Run Metrics HTTP server as non-blocking, goroutine
	go func() {
		metricsServer := NewMetricsServer(fmt.Sprintf("%s:%d", "0.0.0.0", metricsPort), sugar)
		sugar.Infow("Starting metrics http server")
		if err := metricsServer.Server.ListenAndServe(); err != nil {
			sugar.Error(err.Error())
		}
	}()

	//
	// Signal Handling for Graceful exit
	///
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	sugar.Info("Shutting down")
	os.Exit(0)
}
