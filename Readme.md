## Youtrack alerts on new Github release 
A simple tool receives alerts from alertmanager and push a new task on a Youtrack board. That's all.

### How to run
To run it:
```bash
./main --config /config/path
```

### Alertmanager Webhook
In alertmanager config file write:
```yml
receivers:
  - name: "api-webhook"
    webhook_configs:
      - url: "http://metrics-catcher:12000/on_alert"
```
Where `metrics-catcher:12000` are host and port of the tool

### Additional
In addition on `debug` port there is a `/metrics` endpoint for Prometheus metrics