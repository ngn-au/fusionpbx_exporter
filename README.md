# fusionpbx_exporter
Prometheus exporter for FusionPBX Multi-Tenant based metrics

`./fusionpbx_exporter --password YourSecretPassword`

You can find your password in `/etc/fusionpbx/config.php`


Prometheus Config:

```  - job_name: 'fusionpbx'
    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s
    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.
    static_configs:
      - targets: ['fusionpbx:8080']
```