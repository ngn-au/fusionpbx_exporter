package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
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
)

func init() {
	reg.MustRegister(domainCountGauge)
}

func CollectMetrics() {

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
	for _, g := range durationMetrics {
		g.Set(0)
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", *user, *password, *dbname, *host, *port)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	var domainCount sql.NullFloat64
	err = db.QueryRow("SELECT COUNT(*) FROM v_domains").Scan(&domainCount)
	if err != nil {
		log.Fatal(err)
	}

	domainCountGauge.Set(domainCount.Float64)

	rows, err := db.Query("SELECT d.domain_name, COUNT(e.extension) FROM v_domains d JOIN v_extensions e ON d.domain_uuid = e.domain_uuid GROUP BY d.domain_name")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var domain string
		var count sql.NullFloat64
		if err := rows.Scan(&domain, &count); err != nil {
			log.Fatal(err)
		}

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

		gauge.Set(count.Float64)
	}

	rows, err = db.Query("SELECT d.domain_name, COUNT(*) FROM v_domains d JOIN v_xml_cdr c ON d.domain_uuid = c.domain_uuid WHERE c.hangup_cause = 'NORMAL_CLEARING' AND c.end_stamp > clock_timestamp() - INTERVAL '30 seconds' GROUP BY d.domain_name")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var domain string
		var answeredCallsCount sql.NullFloat64
		if err := rows.Scan(&domain, &answeredCallsCount); err != nil {
			log.Fatal(err)
		}

		gauge, ok := answeredCalls[domain]
		if !ok {
			gauge = prometheus.NewGauge(prometheus.GaugeOpts{
				Name:        "answered_calls_per_domain",
				Help:        "Number of answered calls per domain",
				ConstLabels: prometheus.Labels{"domain": domain},
			})
			reg.MustRegister(gauge)
			answeredCalls[domain] = gauge
		}

		gauge.Set(answeredCallsCount.Float64)
	}

	rows, err = db.Query("SELECT domain_name, COUNT(*) FROM v_xml_cdr WHERE direction = 'outbound' AND end_stamp > clock_timestamp() - INTERVAL '30 seconds' GROUP BY domain_name")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var domain string
		var count sql.NullFloat64
		if err := rows.Scan(&domain, &count); err != nil {
			log.Fatal(err)
		}

		gauge, ok := outboundCalls[domain]
		if !ok {
			gauge = prometheus.NewGauge(prometheus.GaugeOpts{
				Name:        "outbound_calls_per_domain",
				Help:        "Number of outbound calls per domain",
				ConstLabels: prometheus.Labels{"domain": domain},
			})
			reg.MustRegister(gauge)
			outboundCalls[domain] = gauge
		}

		gauge.Set(count.Float64)
	}

	rows, err = db.Query("SELECT domain_name, COUNT(*) FROM v_xml_cdr WHERE direction = 'inbound' AND end_stamp > clock_timestamp() - INTERVAL '30 seconds' AND hangup_cause != 'LOSE_RACE' GROUP BY domain_name")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var domain string
		var count sql.NullFloat64
		if err := rows.Scan(&domain, &count); err != nil {
			log.Fatal(err)
		}

		gauge, ok := inboundCalls[domain]
		if !ok {
			gauge = prometheus.NewGauge(prometheus.GaugeOpts{
				Name:        "inbound_calls_per_domain",
				Help:        "Number of inbound calls per domain",
				ConstLabels: prometheus.Labels{"domain": domain},
			})
			reg.MustRegister(gauge)
			inboundCalls[domain] = gauge
		}

		gauge.Set(count.Float64)
	}

	rows, err = db.Query("SELECT domain_name, AVG(rtp_audio_in_mos) FROM v_xml_cdr WHERE end_stamp > clock_timestamp() - INTERVAL '30 seconds' GROUP BY domain_name")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var domain string
		var avgMOS sql.NullFloat64
		if err := rows.Scan(&domain, &avgMOS); err != nil {
			log.Fatal(err)
		}
	
		gauge, ok := mosMetrics[domain]
		if !ok {
			gauge = prometheus.NewGauge(prometheus.GaugeOpts{
				Name:        "average_mos_per_domain",
				Help:        "Average MOS per domain",
				ConstLabels: prometheus.Labels{"domain": domain},
			})
			reg.MustRegister(gauge)
			mosMetrics[domain] = gauge
		}
	
		gauge.Set(avgMOS.Float64)
	}

	rows, err = db.Query("SELECT domain_name, AVG(duration) FROM v_xml_cdr WHERE end_stamp > clock_timestamp() - INTERVAL '30 seconds' GROUP BY domain_name")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var domain string
		var avgDuration sql.NullFloat64
		if err := rows.Scan(&domain, &avgDuration); err != nil {
			log.Fatal(err)
		}
	
		gauge, ok := durationMetrics[domain]
		if !ok {
			gauge = prometheus.NewGauge(prometheus.GaugeOpts{
				Name:        "avg_call_duration_per_domain",
				Help:        "Average call duration per domain in the last 15s",
				ConstLabels: prometheus.Labels{"domain": domain},
			})
			reg.MustRegister(gauge)
			durationMetrics[domain] = gauge
		}
	
		gauge.Set(avgDuration.Float64)
	}
	
	rows.Close()	
    defer db.Close() 

}
