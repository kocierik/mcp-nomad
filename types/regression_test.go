package types

import (
	"encoding/json"
	"testing"
)

// Node.Drivers is an object map in Nomad's API, not map[string]bool.
// Regression for: "json: cannot unmarshal object into Go struct field Node.drivers of type bool".
func TestNodeDecodesRealAPIShape(t *testing.T) {
	const payload = `{
		"ID": "2a6f509a",
		"Name": "client-01",
		"Status": "ready",
		"NodeClass": "batch",
		"NodePool": "default",
		"Resources": {"CPU": 4000, "MemoryMB": 8192, "DiskMB": 51200},
		"Drivers": {
			"docker": {"Detected": true, "Healthy": true, "HealthDescription": "Healthy"},
			"exec":   {"Detected": false, "Healthy": false, "HealthDescription": "Driver not detected"}
		}
	}`
	var n Node
	if err := json.Unmarshal([]byte(payload), &n); err != nil {
		t.Fatalf("decode node: %v", err)
	}
	if !n.Drivers["docker"].Healthy || n.Drivers["exec"].Detected {
		t.Fatalf("unexpected drivers: %+v", n.Drivers)
	}
	if n.NodeClass != "batch" || n.NodePool != "default" {
		t.Fatalf("class/pool not decoded: %q/%q", n.NodeClass, n.NodePool)
	}
	if n.Resources.CPU != 4000 || n.Resources.MemoryMB != 8192 || n.Resources.DiskMB != 51200 {
		t.Fatalf("resources not decoded: %+v", n.Resources)
	}
}

func TestNodeSummaryDecodesNodeClass(t *testing.T) {
	var n NodeSummary
	if err := json.Unmarshal([]byte(`{"ID":"a","Name":"b","Status":"ready","Datacenter":"dc1","NodeClass":"batch","NodePool":"default"}`), &n); err != nil {
		t.Fatal(err)
	}
	if n.NodeClass != "batch" || n.NodePool != "default" {
		t.Fatalf("got %+v", n)
	}
}
