package manifest

import (
	"encoding/json"
	"testing"
)

func TestWorkloadOptionsJSONRoundTrip(t *testing.T) {
	replicas := int32(1)
	want := AppConfiguration{
		WorkloadOptions: WorkloadOptions{
			"frigate": {
				Replicas: &replicas,
				Allow: []WorkloadCapability{
					{Type: WorkloadAllowFolder, Containers: []string{"frigate"}},
					{Type: WorkloadAllowDeviceVideo, Containers: []string{"frigate"}},
				},
				OverlayGateway: &WorkloadOverlayGateway{
					Entrances: []WorkloadOverlayEntrance{{
						Title: "Frigate", Port: 5000, Protocol: "tcp",
					}},
				},
			},
		},
	}

	b, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var got AppConfiguration
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	option := got.WorkloadOptions["frigate"]
	if option.Replicas == nil || *option.Replicas != 1 {
		t.Fatalf("replicas = %v, want 1", option.Replicas)
	}
	if len(option.Allow) != 2 || option.Allow[1].Type != WorkloadAllowDeviceVideo {
		t.Fatalf("allow = %#v", option.Allow)
	}
	if option.OverlayGateway == nil || len(option.OverlayGateway.Entrances) != 1 {
		t.Fatalf("overlayGateway = %#v", option.OverlayGateway)
	}
}
