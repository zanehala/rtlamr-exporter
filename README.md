# rtlamr-exporter

A simple program that exposes [rtlamr](https://github.com/bemasher/rtlamr) messages as [Prometheus](https://prometheus.io/) metrics.

## Usage
Clone this repository and run:
```console
$ rtlamr -format=json | go run .
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

While these messages often contain more information, like tamper flags and backflow metrics, to keep cardinality somewhat low this only exposes three metric labels.
Those labels are:

`protocolType`: The protocol of the message (SCM, SCM+, IDM, NetIDM, R900, R900BCD).

`meterId`: The unique ID of the consumption meter.

`meterType`: Equivalent to ERT type, often known as commodity type (gas, water, power).

### Tips and Tricks
By default rtlamr only listens for SCM messages, you can specify more message types by running it with the `-msgtype` flag:
```console
$ rtlamr -format=json -msgtype=all | rtlamr-exporter
```

Be aware that depending on your location and radio/antenna configuration you may end up with hundreds or thousands of meters, and subsequently time-series.

You can filter to a specific meter ID with the `-filterid` flag:
```console 
$ rtlamr -format=json -filterid=12345 | rtlamr-exporter
```

[See the rtlamr documenting for more.](https://github.com/bemasher/rtlamr/wiki/Configuration#command-line-flags)