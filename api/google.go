package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type GoogleFeed struct {
	URL                   string
	Logger                *zap.SugaredLogger
	Chan                  chan Outage
	PollDurationHistogram *prometheus.HistogramVec
}

func (g *GoogleFeed) getFeed() ([]byte, error) {
	resp, err := http.Get(g.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (g *GoogleFeed) ParseFeed() ([]GoogleOutage, error) {
	status, err := g.getFeed()
	if err != nil {
		return nil, err
	}

	var outages []GoogleOutage

	err = json.Unmarshal(status, &outages)
	if err != nil {
		return nil, err
	}

	return outages, nil
}

func (gf *GoogleFeed) GetOutages() error {
	timer := prometheus.NewTimer(gf.PollDurationHistogram.WithLabelValues(gf.URL))
	defer timer.ObserveDuration()
	googleOutages, err := gf.ParseFeed()
	if err != nil {
		return err
	}

	yesterday := time.Now().AddDate(0, 0, -1)
	for _, outage := range googleOutages {
		if outage.Begin.Before(yesterday) {
			continue
		}

		if (outage.StatusImpact != "SERVICE_DISRUPTION") && (outage.StatusImpact != "SERVICE_OUTAGE") {
			continue
		}

		if outage.Severity != "high" {
			continue
		}

		gf.Chan <- outage.ToOutage()
	}

	return nil
}

func (gf *GoogleFeed) Poll(delay time.Duration) {
	for {
		err := gf.GetOutages()
		if err != nil {
			gf.Logger.Error(err.Error())
		}

		time.Sleep(delay)
	}
}

type GoogleProduct struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type GoogleLocation struct {
	ID    string `json:"ID"`
	Title string `json:'title"`
}

type GoogleOutageUpdate struct {
	Created           time.Time
	Modified          time.Time
	When              time.Time
	Text              string
	Status            string
	AffectedLocations []GoogleLocation `json:"affected_locations"`
}

type GoogleOutage struct {
	ID                          string
	Number                      string
	Begin                       time.Time
	Created                     time.Time
	Modified                    time.Time
	End                         time.Time
	ExternalDescription         string `json:"external_desc"`
	Updates                     []GoogleOutageUpdate
	MostRecentUpdate            GoogleOutageUpdate `json:"most_recent_update"`
	StatusImpact                string             `json:"status_impact"`
	Severity                    string
	ServiceKey                  string           `json:"service_key"`
	ServiceName                 string           `json:"service_name"`
	AffectedProducts            []GoogleProduct  `json:"affected_products"`
	CurrentlyAffectedLocations  []GoogleLocation `json:"currently_affected_locations"`
	PreviouslyAffectedLocations []GoogleLocation `json:"previously_affected_locations"`
}

func (out *GoogleOutage) ToOutage() Outage {
	return Outage{
		Provider:     "google",
		Service:      out.ServiceName,
		StartTime:    out.Begin,
		ModifiedTime: out.Modified,
		Description:  out.ExternalDescription,
	}
}
