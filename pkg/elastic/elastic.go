package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
)

// https://github.com/elastic/go-elasticsearch#go-elasticsearch
func Search(elasticAddr, podName string, logsHits int) (string, error) {
	fmt.Println(podName)
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
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()
		if res.IsError() {
			var e map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				log.Fatalf("Error parsing the response body: %s", err)
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
			log.Fatalf("Error parsing the response body: %s", err)
		}
		if int(r["hits"].(map[string]interface{})["total"].(float64)) <= logsHits {
			fmt.Println("logs lower than hits")
			time.Sleep(2 * time.Second)
		} else {
			logsMatch = int(r["hits"].(map[string]interface{})["total"].(float64))
		}
		//Print the timestamp for each hit.
		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			log.Printf("timestamp=%s", hit.(map[string]interface{})["_source"].(map[string]interface{})["@timestamp"])
		}
		// Print the response status, number of results, and request duration.
		log.Printf(
			"%d hits; took: %dms\n",
			int(r["hits"].(map[string]interface{})["total"].(float64)),
			int(r["took"].(float64)),
		)
		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
		}
		fmt.Println(int(r["hits"].(map[string]interface{})["total"].(float64)))
	}
	return "", nil
}
