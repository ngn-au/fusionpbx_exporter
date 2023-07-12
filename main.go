package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/alecthomas/kingpin/v2"
)

var (
	domainCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "domains_count",
		Help: "Number of domains in v_domains table",
	})
	extensionsCount = make(map[string]prometheus.Gauge)
	answeredCalls   = make(map[string]prometheus.Gauge)
	outboundCalls = make(map[string]prometheus.Gauge)
	inboundCalls = make(map[string]prometheus.Gauge)

	mosMetrics = make(map[string]prometheus.Gauge)
	durationMetrics = make(map[string]prometheus.Gauge)
	
	user     = kingpin.Flag("user", "PostgreSQL username").Default("fusionpbx").String()
	password = kingpin.Flag("password", "PostgreSQL password").Default("password").String()
	dbname   = kingpin.Flag("dbname", "PostgreSQL database name").Default("fusionpbx").String()
	host     = kingpin.Flag("host", "PostgreSQL host").Default("localhost").String()
	port     = kingpin.Flag("port", "PostgreSQL port").Default("5432").String()
)

var reg = prometheus.NewPedanticRegistry()

func init() {
	reg.MustRegister(domainCountGauge)
}

func collectMetrics() {

	for _, g := range extensionsCount {
		g.Set(0)
	}
	for _, g := range answeredCalls {
		g.Set(0)
	}
	for _, g := range outboundCalls {
		g.Set(0)
	}
	for _, g := range inboundCalls {
		g.Set(0)
	}
	for _, g := range mosMetrics {
		g.Set(0)
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", *user, *password, *dbname, *host, *port)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Query the count of domains
	var domainCount float64
	err = db.QueryRow("SELECT COUNT(*) FROM v_domains").Scan(&domainCount)
	if err != nil {
		log.Fatal(err)
	}

	// Set the value of the domain count gauge
	domainCountGauge.Set(domainCount)

	// Execute your query
	rows, err := db.Query("SELECT d.domain_name, COUNT(e.extension) FROM v_domains d JOIN v_extensions e ON d.domain_uuid = e.domain_uuid GROUP BY d.domain_name")
	if err != nil {
		log.Fatal(err)
	}

	// Collect the results
	for rows.Next() {
		var domain string
		var count float64
		if err := rows.Scan(&domain, &count); err != nil {
			log.Fatal(err)
		}

		// Get or create the gauge for the domain
		gauge, ok := extensionsCount[domain]
		if !ok {
			gauge = prometheus.NewGauge(prometheus.GaugeOpts{
				Name:        "extensions_per_domain",
				Help:        "Number of extensions per domain",
				ConstLabels: prometheus.Labels{"domain": domain},
			})
			reg.MustRegister(gauge)
			extensionsCount[domain] = gauge
		}

		// Set the value of the gauge to the result of the query
		gauge.Set(count)
	}

	// Execute your query
	rows, err = db.Query("SELECT d.domain_name, COUNT(*) FROM v_domains d JOIN v_xml_cdr c ON d.domain_uuid = c.domain_uuid WHERE c.hangup_cause = 'NORMAL_CLEARING' AND c.start_stamp > clock_timestamp() - INTERVAL '15 seconds' GROUP BY d.domain_name")
	if err != nil {
		log.Fatal(err)
	}

	// Collect the results
	for rows.Next() {
		var domain string
		var answeredCallsCount float64
		if err := rows.Scan(&domain, &answeredCallsCount); err != nil {
			log.Fatal(err)
		}

		// Get or create the gauge for the domain
		gauge, ok := answeredCalls[domain]
		if !ok {
			gauge = prometheus.NewGauge(prometheus.GaugeOpts{
				Name:        "answered_calls_per_domain_last_15s",
				Help:        "Number of answered calls per domain in the last 15 seconds",
				ConstLabels: prometheus.Labels{"domain": domain},
			})
			reg.MustRegister(gauge)
			answeredCalls[domain] = gauge
		}

		// Set the value of the gauge to the result of the query
		gauge.Set(answeredCallsCount)
	}

	// Execute your query for outbound calls in last 15 seconds
	rows, err = db.Query("SELECT domain_name, COUNT(*) FROM v_xml_cdr WHERE direction = 'outbound' AND start_stamp > clock_timestamp() - INTERVAL '15 seconds' GROUP BY domain_name")
	if err != nil {
		log.Fatal(err)
	}
	// Collect the results
	for rows.Next() {
		var domain string
		var count float64
		if err := rows.Scan(&domain, &count); err != nil {
			log.Fatal(err)
		}

		// Get or create the gauge for the domain
		gauge, ok := outboundCalls[domain]
		if !ok {
			gauge = prometheus.NewGauge(prometheus.GaugeOpts{
				Name:        "outbound_calls_per_domain_last_15s",
				Help:        "Number of outbound calls per domain in the last 15 seconds",
				ConstLabels: prometheus.Labels{"domain": domain},
			})
			reg.MustRegister(gauge)
			outboundCalls[domain] = gauge
		}

		// Set the value of the gauge to the result of the query
		gauge.Set(count)
	}

	// Execute your query for inbound calls in last 15 seconds
	rows, err = db.Query("SELECT domain_name, COUNT(*) FROM v_xml_cdr WHERE direction = 'inbound' AND start_stamp > clock_timestamp() - INTERVAL '15 seconds' GROUP BY domain_name")
	if err != nil {
		log.Fatal(err)
	}
	// Collect the results
	for rows.Next() {
		var domain string
		var count float64
		if err := rows.Scan(&domain, &count); err != nil {
			log.Fatal(err)
		}

		// Get or create the gauge for the domain
		gauge, ok := inboundCalls[domain]
		if !ok {
			gauge = prometheus.NewGauge(prometheus.GaugeOpts{
				Name:        "inbound_calls_per_domain_last_15s",
				Help:        "Number of inbound calls per domain in the last 15 seconds",
				ConstLabels: prometheus.Labels{"domain": domain},
			})
			reg.MustRegister(gauge)
			inboundCalls[domain] = gauge
		}

		// Set the value of the gauge to the result of the query
		gauge.Set(count)
	}

	rows, err = db.Query("SELECT domain_name, AVG(rtp_audio_in_mos) FROM v_xml_cdr WHERE start_stamp > clock_timestamp() - INTERVAL '15 seconds' GROUP BY domain_name")
	if err != nil {
		log.Fatal(err)
	}

	// Collect the results
	for rows.Next() {
		var domain string
		var avgMOS float64
		if err := rows.Scan(&domain, &avgMOS); err != nil {
			log.Fatal(err)
		}
	
		// Get or create the gauge for the domain
		gauge, ok := mosMetrics[domain]
		if !ok {
			gauge = prometheus.NewGauge(prometheus.GaugeOpts{
				Name:        "mos_per_domain_last_15s",
				Help:        "Average MOS per domain in the last 15 seconds",
				ConstLabels: prometheus.Labels{"domain": domain},
			})
			reg.MustRegister(gauge)
			mosMetrics[domain] = gauge
		}
	
		// Set the value of the gauge to the result of the query
		gauge.Set(avgMOS)
	}

	rows, err = db.Query("SELECT domain_name, AVG(duration) FROM v_xml_cdr WHERE start_stamp > clock_timestamp() - INTERVAL '15 seconds' GROUP BY domain_name")
	if err != nil {
		log.Fatal(err)
	}
	// Collect the results
	for rows.Next() {
		var domain string
		var avgDuration float64
		if err := rows.Scan(&domain, &avgDuration); err != nil {
			log.Fatal(err)
		}
	
		// Get or create the gauge for the domain
		gauge, ok := durationMetrics[domain]
		if !ok {
			gauge = prometheus.NewGauge(prometheus.GaugeOpts{
				Name:        "avg_call_duration_per_domain_last_15s",
				Help:        "Average call duration per domain in the last 15s",
				ConstLabels: prometheus.Labels{"domain": domain},
			})
			reg.MustRegister(gauge)
			durationMetrics[domain] = gauge
		}
	
		// Set the value of the gauge to the result of the query
		gauge.Set(avgDuration)
	}
	rows.Close()
	
	
    defer db.Close() // Close the database connection when the function ends or returns
	rows.Close()  // always close rows after you're done with them

}

func main() {
	kingpin.Parse()

	// Collect metrics in a separate goroutine
	go func() {
		for {
			collectMetrics()
			time.Sleep(10 * time.Second)
		}
	}()

	// Expose the metrics for Prometheus to scrape
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(":8080", nil)
}
