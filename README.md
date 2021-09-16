# Prometheus Exporter JUnit

Takes a JUnit XML file, parses it, then writes the metrics to a Prometheus push gateway.
# Install

Option 1: Install with go -> `go get github.com/ianks/prometheus-exporter-junit`
Option 2: Download the tarball from the releases page 
# Usage

```sh
$ prometheus-exporter-junit  -f junit.xml -pushurl  'https://test:t@webhook.site/330677d3-e144-42e2-a529-dd0dc344cee1'
```
