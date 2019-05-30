# beacon
Tiny system beacon reporting CPU, memory, uptime to logs, cloudwatch, http

## Usage

Metrics can be sent to logs (stdout), cloudwatch, or posted as json to an http endpoint.
The output is silent for cloudwatch and http, unless an error occurs.

Here, we send the metrics to both logs, cloudwatch, and http endpoint:

```
# config.yaml
interval: 60s
targets:
  - type: log
  - type: cloudwatch
  - type: url
    arg: http://example.com/heartbeat
```

```
$ beacon
2019/05/30 08:18:05 {Hostname:desktop HostID:83124810-114D-9785-2296-468d43bbbcd5 CPUPercent:8.359375 MemPercent:40 Uptime:6086447}
2019/05/30 08:19:05 {Hostname:desktop HostID:83124810-114D-9785-2296-468d43bbbcd5 CPUPercent:7.682194 MemPercent:40 Uptime:6086507}
...
```

## License

[MIT License](https://github.com/shazbits/beacon/blob/master/LICENSE)
