package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAmazonFeedsGetFeeds(t *testing.T) {
	testRegion := "test-region"
	testProduct := "test-product"

	azFeeds := AmazonFeeds{
		Regions:  []string{testRegion},
		Products: []string{testProduct},
	}

	feeds := azFeeds.GetFeeds()

	if len(feeds) == 0 {
		t.Fatalf(`AmazonFeeds length is zero`)
	}

	expectedUrl := fmt.Sprintf("https://status.aws.amazon.com/rss/%s-%s.rss", testProduct, testRegion)
	actualUrl := azFeeds.GetFeeds()[0].URL
	if actualUrl != expectedUrl {
		t.Fatalf(`URL was %s, expected %s`, actualUrl, expectedUrl)
	}
}

func TestAmazonFeedGetOutagesIgnoreTitleString(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `
		<?xml version="1.0" encoding="UTF-8"?>
		<rss version="2.0">
		  <channel>
			<title><![CDATA[Amazon Elastic Container Service (N. California) Service Status]]></title>
			<link>http://status.aws.amazon.com/</link>
			<language>en-us</language>
			<lastBuildDate>Mon, 29 Aug 2022 17:03:35 PDT</lastBuildDate>
			<generator>AWS Service Health Dashboard RSS Generator</generator>
			<description><![CDATA[Amazon Elastic Container Service (N. California) Service Status]]></description>
			<ttl>5</ttl>
			<!-- You seem to care about knowing about your events, why not check out https://docs.aws.amazon.com/health/latest/ug/getting-started-api.html -->
		
			 
			 <item>
			  <title><![CDATA[Service is operating normally: [RESOLVED] Elevated API Error Rates]]></title>
			  <link>http://status.aws.amazon.com/</link>
			  <pubDate>Sat, 09 May 2020 21:52:26 PDT</pubDate>
			  <guid isPermaLink="false">http://status.aws.amazon.com/#ecs-us-west-1_1589086346</guid>
			  <description><![CDATA[Between 6:53 PM and 9:37 PM PDT we experienced elevated API error rates and latencies in the US-WEST-1 Region. The issue has been resolved and the service is operating normally. Running tasks were not impacted.]]></description>
			 </item>
				
		  </channel>
		</rss>
		`)
	}))
	defer ts.Close()

	logger := zaptest.NewLogger(t)
	azFeed := AmazonFeed{
		URL:    ts.URL,
		Logger: logger.Sugar(),
		PollDurationHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "C",
			Help: "Help",
		}, []string{"handler"}),
	}

	outages, err := azFeed.GetOutages()
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Equal(t, len(outages), 0, "There should be zero outages reported")
}

func TestAmazonFeedGetOutages(t *testing.T) {
	fmtDate := time.Now().Format(time.RFC1123)
	fmtResponse := fmt.Sprintf(`
	<?xml version="1.0" encoding="UTF-8"?>
	<rss version="2.0">
	  <channel>
		<title><![CDATA[Amazon Elastic Container Service (N. California) Service Status]]></title>
		<link>http://status.aws.amazon.com/</link>
		<language>en-us</language>
		<lastBuildDate>%s</lastBuildDate>
		<generator>AWS Service Health Dashboard RSS Generator</generator>
		<description><![CDATA[Amazon Elastic Container Service (N. California) Service Status]]></description>
		<ttl>5</ttl>		
		 
		 <item>
		  <title><![CDATA[Service is operating normally: Elevated API Error Rates]]></title>
		  <link>http://status.aws.amazon.com/</link>
		  <pubDate>%s</pubDate>
		  <guid isPermaLink="false">http://status.aws.amazon.com/#ecs-us-west-1_1589086346</guid>
		  <description><![CDATA[Between 6:53 PM and 9:37 PM PDT we experienced elevated API error rates and latencies in the US-WEST-1 Region. The issue has been resolved and the service is operating normally. Running tasks were not impacted.]]></description>
		 </item>
			
	  </channel>
	</rss>
	`, fmtDate, fmtDate)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmtResponse)
	}))
	defer ts.Close()

	logger := zaptest.NewLogger(t)
	azFeed := AmazonFeed{
		URL:    ts.URL,
		Logger: logger.Sugar(),
		PollDurationHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "C",
			Help: "Help",
		}, []string{"handler"}),
	}

	outages, err := azFeed.GetOutages()
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Equal(t, len(outages), 1, "There should be one outage reported")
}
