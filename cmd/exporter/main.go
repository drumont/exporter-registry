package main

import (
	"bytes"
	"flag"
	"fmt"
	exporter_registry "github.com/drumont/exporter-registry"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {

	var (
		targetHost = flag.String("target.host", "localhost", "nginx basic status page")
		targetPort = flag.Int("target.port", 8080, "nginx basic status page")
		targetPath = flag.String("target.path", "/status", "nginx status page")
		promPort   = flag.Int("prom.port", 9150, "Port to expose metrics")
	)
	flag.Parse()

	uri := fmt.Sprintf("http://%s:%d%s", *targetHost, *targetPort, *targetPath)

	//Called on each collector collect
	basicStats := func() ([]exporter_registry.NginxStats, error) {
		var netClient = &http.Client{
			Timeout: time.Second * 10,
		}
		resp, err := netClient.Get(uri)
		if err != nil {
			log.Fatalf("Error getting nginx stats: %v", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading nginx stats: %v", err)
		}
		r := bytes.NewReader(body)
		return exporter_registry.ScanBasicStats(r)
	}

	bc := exporter_registry.NewBasicCollector(basicStats)

	reg := prometheus.NewRegistry()
	reg.MustRegister(bc)

	mux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", promHandler)

	port := fmt.Sprintf(":%d", *promPort)
	log.Printf("Starting server on port %s", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
