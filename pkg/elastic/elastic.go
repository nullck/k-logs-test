package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/nullck/k-logs-test/pkg/prometheus_push"
)

type ES struct {
	ElasticAddr string
	PodName     string
	LogsHits    int
	Threshold   int
}

var buf bytes.Buffer
var r map[string]interface{}
var ej map[string]interface{}
var timeLayout = "2006-01-02T15:04:05"
var status = "OK"

func (e *ES) getIndex() (elasticAddr, indexName string) {
	i := strings.Split(e.ElasticAddr, "/")
	indexName = i[3]
	elasticAddr = strings.Replace(e.ElasticAddr, "/"+indexName, "", 1)
	return elasticAddr, indexName
}

func (e *ES) getTotal() int {
	elasticAddr, indexName := e.getIndex()
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticAddr,
		},
	}
	countQuery := fmt.Sprintf("kubernetes.pod_name: \"%s\"", string(e.PodName))
	es, _ := elasticsearch.NewClient(cfg)
	res, err := es.Count(
		es.Count.WithContext(context.Background()),
		es.Count.WithIndex(indexName),
		es.Count.WithQuery(string(countQuery)),
		es.Count.WithPretty(),
	)
	if err != nil {
		log.Fatalf("error getting response from getTotal func: %s", err)
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("error parsing the response body: %s", err)
	}

	return int(r["count"].(float64))
}

func (e *ES) Search(promEnabled bool, promGWAddr string, promGWPort int) (string, error) {
	elasticAddr, indexName := e.getIndex()
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticAddr,
		},
	}
	es, _ := elasticsearch.NewClient(cfg)
	logsMatch := 0

	for logsMatch < e.LogsHits {
		query := fmt.Sprintf("{\"query\": { \"match\": { \"kubernetes.pod_name\": {\"query\": \"%s\", \"operator\": \"and\"}}}}", e.PodName)
		res, err := es.Search(
			es.Search.WithContext(context.Background()),
			es.Search.WithIndex(indexName),
			es.Search.WithTrackTotalHits(true),
			es.Search.WithPretty(),
			es.Search.WithErrorTrace(),
			es.Search.WithBody(
				strings.NewReader(query)),
		)
		logsMatch = e.getTotal()

		if err != nil {
			log.Fatalf("error getting response: %s", err)
		}
		defer res.Body.Close()
		if res.IsError() {
			if err := json.NewDecoder(res.Body).Decode(&ej); err != nil {
				log.Fatalf("error parsing the response body: %s", err)
			} else {
				// Print the response status and error information.
				log.Fatalf("[%s] %s: %s",
					res.Status(),
					ej["error"].(map[string]interface{})["type"],
					ej["error"].(map[string]interface{})["reason"],
				)
			}
		}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Fatalf("error parsing the response body: %s", err)
		}

		if logsMatch < e.LogsHits {
			log.Printf("total logs lower than log-hits specified ... wait")
			time.Sleep(500 * time.Millisecond)
		}

		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			elasticTimestamp := fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["@timestamp"])
			elasticTime := strings.Split(elasticTimestamp, ".")

			containerMsgTimestamp := fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["log"])
			containerTime := regexp.MustCompile(`\d{4}\-\d{1,2}\-\d{1,2}T\d{1,2}\:\d{1,2}\:\d{1,2}$`).FindAllString(containerMsgTimestamp, 1)

			elasticTimeP, _ := time.Parse(timeLayout, elasticTime[0])
			containerTimeP, _ := time.Parse(timeLayout, containerTime[0])
			timeDiff := elasticTimeP.Sub(containerTimeP).Milliseconds()

			log.Printf("container timestamp=%s\n elastic timestamp=%s", containerTimeP, elasticTimeP)

			if promEnabled {
				promMetric(timeDiff, promGWAddr, promGWPort)
			}

			if e.Threshold > 0 {
				if float64(e.Threshold) < float64(timeDiff) {
					status = "ALERT"
				}
			}
			log.Printf("logs delayed in: %v milliseconds", timeDiff)
		}
		log.Printf("total logs %d", logsMatch)
	}
	return status, nil
}

func promMetric(timeDiff int64, promGWAddr string, promGWPort int) {
	var prom = prometheus_push.PrometheusPusher{
		GWUrl:      promGWAddr,
		GWPort:     promGWPort,
		MetricName: "k-logs-delay",
	}
	prom.PushMetric(timeDiff)
}
