package elastic

import "testing"

func TestSearch(t *testing.T) {
	e := ES{"http://localhost:9200/test_logs", "pod-tt", 0, 1}
	s, err := e.Search()
	if s != "OK" {
		t.Errorf("error %v", s)
	}
	if err != nil {
		t.Errorf("error %v", err)
	}
}
