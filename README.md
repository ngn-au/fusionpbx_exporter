# fusionpbx_exporter

    Number of extensions per domain: For each domain in the v_domains table, the number of associated extensions in the v_extensions table is counted and set in the corresponding gauge in the extensionsCount map.

    Number of answered calls per domain: For each domain in the v_domains table, the number of answered calls in the v_xml_cdr table (with a hangup cause of 'NORMAL_CLEARING' in the last 30 seconds) is counted and set in the corresponding gauge in the answeredCalls map.

    Number of missed calls per domain: Similarly, the number of missed calls per domain (with a hangup cause of 'ORIGINATOR_CANCEL' in the last 30 seconds) is counted and set in the missedCalls map.

    Number of outbound calls per domain: The number of outbound calls per domain (in the last 30 seconds) is counted and set in the outboundCalls map.

    Number of inbound calls per domain: The number of inbound calls per domain (that did not have a hangup cause of 'LOSE_RACE' in the last 30 seconds) is counted and set in the inboundCalls map.

    Average MOS per domain: The average Mean Opinion Score (MOS) per domain (calculated over the last 30 seconds) is set in the mosMetrics map.

    Average call duration per domain: The average call duration per domain (calculated over the last 30 seconds) is set in the durationMetrics map.

Each time a new domain is encountered in a query, a new gauge is created with the domain as a constant label, registered to the Prometheus registry, and added to the corresponding map. Existing gauges are reused.

The function ends by closing the rows object and deferring the closure of the database connection. It's good practice to ensure that all open connections and objects that can be closed, like db and rows here, are closed when they are no longer needed. This can prevent resource leaks.

This script will create a number of Prometheus metrics based on data in a FusionPBX database, which can then be scraped by a Prometheus server for monitoring and alerting purposes.

Note that while this script seems to be functioning correctly, proper error handling around database operations and logging would improve its robustness. Moreover, the script seems to connect to the database and fetch the metrics every 10 seconds. This can be resource intensive, especially if the database is large or the queries are complex. One solution to this might be to increase the metrics collection interval, or implement some caching mechanism.



Prometheus exporter for FusionPBX Multi-Tenant based metrics
<img width="1459" alt="Screenshot 2023-07-12 at 3 42 49 pm" src="https://github.com/ngn-au/fusionpbx_exporter/assets/107200645/28feda6d-fcc6-48b0-b6fd-7625b8d48fd4">

Inspired by: https://github.com/florentchauveau/freeswitch_exporter

`./fusionpbx_exporter --password YourSecretPassword`

You can find your password in `/etc/fusionpbx/config.php`



**Prometheus Config:**

```yaml
  - job_name: 'fusionpbx'
    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s
    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.
    static_configs:
      - targets: ['fusionpbx:8080']
```


**Usage:**
```bash
fusionpbx_exporter --help
usage: fusionpbx_exporter [<flags>]


Flags:
  --[no-]help            Show context-sensitive help (also try --help-long and
                         --help-man).
  --user="fusionpbx"     PostgreSQL username
  --password="password"  PostgreSQL password
  --dbname="fusionpbx"   PostgreSQL database name
  --host="localhost"     PostgreSQL host
  --port="5432"          PostgreSQL port

```
**METRICS**
```
# HELP extensions_per_domain Number of extensions per domain
# TYPE extensions_per_domain gauge
# HELP answered_calls_per_domain Number of answered calls per domain
# TYPE answered_calls_per_domain gauge
# HELP missed_calls_per_domain Number of missed calls per domain
# TYPE missed_calls_per_domain gauge
# HELP outbound_calls_per_domain Number of outbound calls per domain
# TYPE outbound_calls_per_domain gauge
# HELP inbound_calls_per_domain Number of inbound calls per domain
# TYPE inbound_calls_per_domain gauge
# HELP avg_mos_per_domain Average MOS per domain
# TYPE avg_mos_per_domain gauge
# HELP avg_call_duration_per_domain Average call duration per domain
# TYPE avg_call_duration_per_domain gauge
```
**Grafana Dashboard**

https://grafana.com/grafana/dashboards/19155-fusionpbx/
