package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/joshdk/go-junit"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/prometheus/common/expfmt"
)

var (
	registry      = prometheus.NewRegistry()
	xmlPath       = flag.String("f", "junit.xml", "Path for the JUnit XML file")
	prometheusUrl = flag.String("pushurl", "", "URL for Prometheus push gateway")

	// testDurationsHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	// 	Name:    "test_duration_seconds",
	// 	Help:    "Test duration histograms.",
	// 	Buckets: prometheus.DefBuckets,
	// },
	// 	[]string{"id", "suite"})
	testCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "testcase_total",
		Help: "A counter of tests that were run.",
	},
		[]string{"id", "suite", "status", "name"})
)

func init() {
	// Register the summary and the histogram with Prometheus's default registry.
	// registry.MustRegister(testDurationsHistogram)
	registry.MustRegister(testCounter)
}

func gatherMetrics() {
	xml, err := ioutil.ReadFile(*xmlPath)
	if err != nil {
		log.Fatalf("failed to ingest JUnit xml %v", err)
	}

	suites, err := junit.Ingest(xml)
	if err != nil {
		log.Fatalf("failed to ingest JUnit xml %v", err)
	}

	for _, suite := range suites {
		fmt.Printf("Analyzing suite %s\n", suite.Name)
		for _, test := range suite.Tests {
			testCounter.WithLabelValues(
				test.Properties["id"],
				suite.Name,
				string(test.Status),
				test.Name,
			).Add(1)
		}
	}
}

func main() {
	flag.Parse()
	gatherMetrics()
	push.New(*prometheusUrl, "ci").Gatherer(registry).Format(expfmt.FmtText).Push()
}
