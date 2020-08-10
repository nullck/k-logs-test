## K-Logs-Test


K-logs-test helps you to be sure that your applications logs are getting into their destination (ElasticSearch) in a proper period of time.

In case the logs are taking more time than your team accept, k-logs-test will identify the problem before and notify you.

#### How it works?

k-logs-test will create a pod in your kubernetes in a namespace you've defined.


#### How to use the k-logs-test utility

Example:

```
k-logs-test run -p tatata -e http://localhost:9200/test_logs --logs-hits 4 --threshold 3 --channel "#k-logs" --webhook-url https://hooks.slack.com/services/XXXX/XXXX/XXXXX --slack-alert-enabled true
```

#### cli reference

```
k-logs-test run --help

Usage:
  k-logs-test run [flags]

Flags:
  -c, --channel string            The Slack Channel for notification (default "#k-logs")
  -e, --elastic-endpoint string   The ElasticSearch Endpoint and the logs index name (default "https://localhost:9200/fluentd")
  -h, --help                      help for run
      --logs-hits int             The number of logs hits (default 30)
  -n, --namespace string          The pod namespace (default "default")
  -p, --pod-name string           The pod name (default "k-logs-test")
  -a, --slack-alert-enabled       Enable or not slack alerts
      --threshold int             The Alert Threshould in milliseconds
  -w, --webhook-url string        The Slack Webhook Url for notification

Global Flags:
      --config string   config file (default is $HOME/.k-logs-test.yaml)
```
