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
)

// https://github.com/elastic/go-elasticsearch#go-elasticsearch
func Search(elasticAddr, podName string, logsHits int) (string, error) {
	log.Printf("podName %s", podName)
	i := strings.Split(elasticAddr, "/")
	indexName := i[3]
	elasticAddr = strings.Replace(elasticAddr, "/"+indexName, "", 1)

	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticAddr,
		},
	}
	es, _ := elasticsearch.NewClient(cfg)

	var buf bytes.Buffer
	var r map[string]interface{}
	var timeLayout = "2006-01-02T15:04:05"
	logsMatch := 0

	for logsMatch <= logsHits {
		query := map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"kubernetes.pod_name": podName,
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
			var e map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				log.Fatalf("error parsing the response body: %s", err)
			} else {
				// Print the response status and error information.
				log.Fatalf("[%s] %s: %s",
					res.Status(),
					e["error"].(map[string]interface{})["type"],
					e["error"].(map[string]interface{})["reason"],
				)
			}
		}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Fatalf("error parsing the response body: %s", err)
		}
		if int(r["hits"].(map[string]interface{})["total"].(float64)) <= logsHits {
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

			log.Printf("container log timestamp=%s", containerTimeP)
			log.Printf("elasticsearch log timestamp=%s", elasticTimeP)

			log.Printf("logs delayed in: %v seconds", elasticTimeP.Sub(containerTimeP).Seconds())

		}
		log.Printf("total logs hits %d", int(r["hits"].(map[string]interface{})["total"].(float64)))
	}
	return "", nil
}
