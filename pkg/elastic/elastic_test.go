package elastic

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

var url = "http://localhost:9200/test_logs/test_logs_type"
var sSeed = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

type TSearch struct {
	podName string
	status  string
}

var promEnabled bool

func (ts *TSearch) GenerateJson() []byte {
	times := "2020-08-02T11:19:57"
	containerLog := "2020-08-02T11:19:20.005Z stdout F tatata: 2020-08-02T11:19:20"

	if ts.status == "OK" {
		times = "2020-08-02T11:19:27"
		containerLog = "2020-08-02T11:19:27 stdout F tatata: 2020-08-02T11:19:27"
	}

	js := fmt.Sprintf(`{
    "@timestamp": "%v",
    "log": "%v",
    "kubernetes": {
      "pod_name": "%v",
      "namespace_name": "default"
    }
  }`, times, containerLog, ts.podName)
	jsonStr := []byte(fmt.Sprintf("%v", js))
	return jsonStr
}

func (ts *TSearch) MockSearch(promEnabled bool, jsonStr []byte) (string, error) {
	fmt.Printf("%v\n", string(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	e := &ES{
		ElasticAddr: "http://localhost:9200/test_logs",
		PodName:     ts.podName,
		LogsHits:    1,
		Threshold:   2000,
	}

	promGWAddr := "prometheus-pushgateway"
	promGWPort := 9090
	s, err := e.Search(promEnabled, promGWAddr, promGWPort)
	return s, err
}

func TestSearch(t *testing.T) {
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(sSeed)
	ts := &TSearch{
		podName: "duck-" + sSeed[n],
		status:  "alert",
	}
	elasticJS := ts.GenerateJson()
	promEnabled = false
	s, err := ts.MockSearch(promEnabled, elasticJS)
	if s != status {
		t.Errorf("error %v", s)
	}
	if err != nil {
		t.Errorf("error %v", err)
	}
}

func TestSearchOK(t *testing.T) {
	n := rand.Int() % len(sSeed)
	ts := &TSearch{
		podName: "chicken-" + sSeed[n],
		status:  "OK",
	}
	elasticJS := ts.GenerateJson()
	promEnabled = true
	s, err := ts.MockSearch(promEnabled, elasticJS)
	if s != status {
		t.Errorf("error %v", s)
	}
	if err != nil {
		t.Errorf("error %v", err)
	}
}
