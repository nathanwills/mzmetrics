package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nathanwills/mzmetrics/pkg/myzone"
	"github.com/nathanwills/mzmetrics/pkg/myzone/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Define command-line flags
	bind := flag.String("bind", ":2112", "binding to serve metrics on")
	acURL := flag.String("ac-url", "http://192.168.1.113:2025/getSystemData", "AC data URL")

	// Parse the flags
	flag.Parse()

	if envBind := os.Getenv("MZMETRICS_BIND"); envBind != "" {
		bind = &envBind
	}

	if envACURL := os.Getenv("MZMETRICS_AC_URL"); envACURL != "" {
		acURL = &envACURL
	}

	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	// Create new metrics and register them using the custom registry.
	m := metrics.New(reg)

	// Create a new context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling (SIGTERM)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)

	// Create a non-blocking timer that ticks every 10 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	ac, err := myzone.Fetch(*acURL)
	if err != nil {
		fmt.Printf("error: fetch: %v\n", err)
		return
	}

	m.SetMetrics(ac)

	go func() {
		for {
			select {
			case <-ctx.Done():
				// Ctx cancelled, exit the routine
				return
			case <-ticker.C:
				ac, err := myzone.Fetch(*acURL)
				if err != nil {
					fmt.Printf("error: fetch: %v\n", err)
					continue
				}

				m.SetMetrics(ac)
			}
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	fmt.Printf("Starting\n")
	http.ListenAndServe(*bind, nil)
}
