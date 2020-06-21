package elastic

import (
	"log"

	"github.com/elastic/go-elasticsearch/v7"
)

// https://github.com/elastic/go-elasticsearch#go-elasticsearch
func main() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, _ := elasticsearch.NewClient(cfg)
	log.Println(elasticsearch.Version)
	log.Println(es.Info())
}
