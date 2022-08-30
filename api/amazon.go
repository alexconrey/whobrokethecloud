package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"strings"
	"time"
)

var (
	ignoreTitleStrings = []string{
		"Impairments",
		"[RESOLVED]",
	}
)

// func getAmazonFeeds() []string {
// 	rtn := []string{}
// 	for _, svc := range []string{
// 		"ec2",
// 		"ecs",
// 		"eks",
// 		"elb",
// 		"elasticache",
// 		"rds",
// 	} {
// 		for _, region := range []string{
// 			"us-east-1",
// 			"us-east-2",
// 			"us-west-1",
// 			"us-west-2",
// 		} {
// 			str := fmt.Sprintf("https://status.aws.amazon.com/rss/%s-%s.rss", svc, region)
// 			rtn = append(rtn, str)
// 		}
// 	}
// 	return rtn
// }

func checkIgnoreTitleStrings(str string) bool {
	for _, ignoreStr := range ignoreTitleStrings {
		if strings.Contains(str, ignoreStr) {
			return true
		}
	}
	return false
}

type AmazonFeed struct {
	URL                   string
	Logger                *zap.SugaredLogger
	PollDurationHistogram *prometheus.HistogramVec
}

func (a *AmazonFeed) GetOutages() ([]Outage, error) {
	timer := prometheus.NewTimer(a.PollDurationHistogram.WithLabelValues(a.URL))
	defer timer.ObserveDuration()
	outages := []Outage{}
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(a.URL)
	if err != nil {
		return nil, err
	}

	yesterday := time.Now().AddDate(0, 0, -31)

	for _, item := range feed.Items {
		if item.PublishedParsed.Before(yesterday) {
			continue
		}

		if checkIgnoreTitleStrings(item.Title) {
			a.Logger.Debugw("Ignoring outage due to title",
				"title", item.Title,
				"provider", "amazon",
				"outage", item,
			)
			continue
		}

		svcName := strings.Trim(feed.Title, "Amazon")
		svcName = strings.Trim(svcName, "Service Status")

		outage := Outage{
			Provider:  "amazon",
			Service:   svcName,
			StartTime: *item.PublishedParsed,
		}

		a.Logger.Infow("Loaded outage information",
			"provider", outage.Provider,
		)

		outages = append(outages, outage)
	}

	return outages, nil
}

type AmazonFeeds struct {
	Regions               []string
	Products              []string
	Chan                  chan Outage
	Logger                *zap.SugaredLogger
	PollDurationHistogram *prometheus.HistogramVec
}

func (az *AmazonFeeds) GetFeeds() []AmazonFeed {
	feeds := []AmazonFeed{}
	for _, svc := range az.Products {
		for _, region := range az.Regions {
			feed := AmazonFeed{
				URL:                   fmt.Sprintf("https://status.aws.amazon.com/rss/%s-%s.rss", svc, region),
				Logger:                az.Logger,
				PollDurationHistogram: az.PollDurationHistogram,
			}
			feeds = append(feeds, feed)
		}
	}
	return feeds
}

func (az *AmazonFeeds) GetOutages() error {
	for _, feed := range az.GetFeeds() {
		feedOutages, err := feed.GetOutages()
		if err != nil {
			return err
		}

		for _, outage := range feedOutages {
			az.Chan <- outage
			az.Logger.Debugw("Added outage event to queue",
				"provider", outage.Provider,
			)
		}
	}
	return nil
}

func (az *AmazonFeeds) Poll(delaySeconds int) error {
	for {
		err := az.GetOutages()
		if err != nil {
			return err
		}

		time.Sleep(time.Second *
			(time.Duration(delaySeconds) * time.Millisecond))
	}
}
