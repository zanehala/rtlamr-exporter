# rtlamr-exporter

A simple program that exposes [rtlamr](https://github.com/bemasher/rtlamr) messages as [Prometheus](https://prometheus.io/) metrics.

## Usage
Clone this repository and run:
```console
$ rtlamr -format=json | go run main.go
```
Or if you've built the application as a binary: 
```console
$ rtlamr -format=json | rtlamr-exporter
```

This should begin serving metrics on `localhost:9090`.
You can view those metrics with:
```console
$ curl <your machines ip or localhost>:9090/metrics
```

Next you'll want to set up a [Prometheus scrape config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config) to scrape these metrics:
```yaml
scrape_configs:
- job_name: rtlamr-exporter
  scrape_interval: 10s # set this to whatever interval you want
  static_configs:
  - targets:
    - localhost:9090 # your machines IP/hostname here
```
