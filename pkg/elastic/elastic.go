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
var logsMatch = 0

func (e *ES) Search(promEnabled bool, promGWAddr string, promGWPort int) (string, error) {
	i := strings.Split(e.ElasticAddr, "/")
	indexName := i[3]
	elasticAddr := strings.Replace(e.ElasticAddr, "/"+indexName, "", 1)

	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticAddr,
		},
	}
	es, _ := elasticsearch.NewClient(cfg)

	for logsMatch < e.LogsHits {
		query := map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"kubernetes.pod_name": e.PodName,
				},
			},
		}
		if err := json.NewEncoder(&buf).Encode(query); err != nil {
			log.Fatalf("Error encoding query: %s", err)
		}
		res, err := es.Search(
			es.Search.WithContext(context.Background()),
			es.Search.WithIndex(indexName),
			es.Search.WithBody(&buf),
			es.Search.WithTrackTotalHits(true),
			es.Search.WithPretty(),
		)
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
		if int(r["hits"].(map[string]interface{})["total"].(float64)) < e.LogsHits {
			log.Printf("total logs lower than log-hits specified ... wait")
			time.Sleep(1 * time.Second)
		} else {
			logsMatch = int(r["hits"].(map[string]interface{})["total"].(float64))
		}

		re := regexp.MustCompile(`\d{4}\-\d{1,2}\-\d{1,2}T\d{1,2}\:\d{1,2}\:\d{1,2}$`)

		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			elasticTimestamp := fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["@timestamp"])
			elasticTime := strings.Split(elasticTimestamp, ".")

			containerMsgTimestamp := fmt.Sprintf("%v", hit.(map[string]interface{})["_source"].(map[string]interface{})["log"])
			containerTime := re.FindAllString(containerMsgTimestamp, 1)

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
		log.Printf("total logs %d", int(r["hits"].(map[string]interface{})["total"].(float64)))
	}
	return status, nil
}

func promMetric(timeDiff int64, promGWAddr string, promGWPort int) {
	type p = prometheus_push.Prometheus
	var prom = p{
		GWUrl:      promGWAddr,
		GWPort:     promGWPort,
		MetricName: "k-logs-delay",
	}
	prom.PushMetric(timeDiff)
}
