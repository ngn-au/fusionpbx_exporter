# fusionpbx_exporter
Prometheus exporter for FusionPBX Multi-Tenant based metrics
<img width="1459" alt="Screenshot 2023-07-12 at 3 42 49 pm" src="https://github.com/ngn-au/fusionpbx_exporter/assets/107200645/28feda6d-fcc6-48b0-b6fd-7625b8d48fd4">

Inspired by: https://github.com/florentchauveau/freeswitch_exporter

`./fusionpbx_exporter --password YourSecretPassword`

You can find your password in `/etc/fusionpbx/config.php`



**Prometheus Config:**

```  - job_name: 'fusionpbx'
    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s
    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.
    static_configs:
      - targets: ['fusionpbx:8080']
```


**Usage:**
```
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
