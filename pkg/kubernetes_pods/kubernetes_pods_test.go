package kubernetes_pods

import (
	"testing"
)

const wantPo = "pod-tt"

func TestCreatePod(t *testing.T) {
	p := Pod{wantPo, "default"}
	logHits := 2
	gotPo, err := p.CreatePod(logHits)
	if gotPo != wantPo {
		t.Errorf("CreatePod wants pod name %v, but got %v", wantPo, gotPo)
	}
	if err != nil {
		t.Errorf("CreatePod err: %v", err)
	}
}

func TestDeletePod(t *testing.T) {
	p := Pod{wantPo, "default"}
	gotPo, err := p.DeletePod(p.PodName)
	if gotPo != wantPo {
		t.Errorf("DeletePod wants pod name %v, but got %v", wantPo, gotPo)
	}
	if err != nil {
		t.Errorf("CreatePod err: %v", err)
	}
}
