# fusionpbx_exporter
Prometheus exporter for FusionPBX Multi-Tenant based metrics

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
