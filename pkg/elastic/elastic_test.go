package elastic

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func init() {
	url := "http://localhost:9200/test_logs/test_logs_type"
	var jsonStr = []byte(`{
    "@timestamp": "2020-08-02T11:19:27.005Z",
    "log": "2020-08-02T11:19:24.1416191Z stdout F tatata: 2020-08-02T11:19:24",
    "kubernetes": {
      "pod_name": "pod_test",
      "namespace_name": "default"
    }
  }`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func TestSearch(t *testing.T) {
	e := ES{
		ElasticAddr: "http://localhost:9200/test_logs",
		PodName:     "pod_test",
		LogsHits:    0,
		Threshold:   1,
	}
	s, err := e.Search()
	if s != "ALERT" {
		t.Errorf("error %v", s)
	}
	if err != nil {
		t.Errorf("error %v", err)
	}
}
