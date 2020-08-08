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



