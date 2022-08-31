package main

import (
	// "fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Outage struct {
	Provider     string
	Service      string
	StartTime    time.Time
	ModifiedTime time.Time
	Description  string
}

type Outages struct {
	mu      sync.Mutex
	Outages []Outage
	Chan    chan Outage
	Logger  *zap.SugaredLogger
}

var (
	outagesProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "outage_events_processed",
		Help: "The total number of outage events processed",
	},
		[]string{"provider"},
	)

	outagesRemoved = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "outage_events_removed",
		Help: "The total number of outage events removed",
	},
		[]string{"provider"},
	)

	outagesGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "outage_events_gauge",
		Help: "The current number of outage events",
	},
	)
)

func (o *Outages) PollMetrics() {
	for {
		outagesGauge.Set(float64(len(o.Outages)))
		time.Sleep(time.Second / 2)
	}
}

func (o *Outages) HandleOutages() {
	o.Logger.Infow("Starting outage handler")
	// for {
	// 	select {
	// 	case outage := <-o.Chan:
	// 		o.Logger.Infow("Processing outage event",
	// 			"provider", outage.Provider,
	// 			"start_time", outage.StartTime,
	// 		)
	// 		o.AddOutage(outage)
	// 		outagesProcessed.With(prometheus.Labels{"provider": outage.Provider}).Inc()
	// 	}
	// }
	for {
		outage := <-o.Chan
		o.Logger.Infow("Processing outage event",
			"provider", outage.Provider,
			"start_time", outage.StartTime,
		)
		o.AddOutage(outage)
		outagesProcessed.With(prometheus.Labels{"provider": outage.Provider}).Inc()
	}
}

func (o *Outages) removeIndexFromOutages(i int) {
	out := []Outage{}
	for idx, item := range o.Outages {
		if idx == i {
			o.Logger.Debugw("Removing index from outages",
				"index", idx,
				"outage", item,
			)
			outagesRemoved.With(prometheus.Labels{"provider": item.Provider}).Inc()
			continue
		}
		out = append(out, item)
	}
	o.mu.Lock()
	defer o.mu.Unlock()

	o.Outages = out
}

func (o *Outages) AddOutages(outages []Outage) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Outages = append(o.Outages, outages...)
}

func (o *Outages) AddOutage(outage Outage) {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, existingOutage := range o.Outages {
		if outage == existingOutage {
			o.Logger.Debugw("Skipping outage event as it already exists",
				"provider", outage.Provider,
				"start_time", outage.StartTime,
			)
			return
		}
	}
	o.Outages = append(o.Outages, outage)
	o.Logger.Debugw("Loaded outage event",
		"provider", outage.Provider,
		"start_time", outage.StartTime,
	)
}

func (o *Outages) StartWatchdog(duration time.Duration) {
	o.Logger.Infow("Starting TTL Watchdog")
	// TTL is passed as a unit of days, therefore we need to invert
	// and look that many negative days ago
	invert_date := duration * -1
	max_date := time.Now().Add(invert_date)

	for {
		for idx, item := range o.Outages {
			// Check if start or modified time exceed the max date
			if item.StartTime.Before(max_date) && item.ModifiedTime.Before(max_date) {
				o.Logger.Infow(
					"Removing due to expired TTL",
					// "outage", item,
				)
				o.removeIndexFromOutages(idx)
			}
		}
		time.Sleep(duration)
	}
}

func NewOutages(logger *zap.SugaredLogger) Outages {
	return Outages{
		Outages: make([]Outage, 0),
		Chan:    make(chan Outage),
		Logger:  logger,
	}
}
