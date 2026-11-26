package manifest

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
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

func TestWorkloadAcceleratorYAML(t *testing.T) {
	src := []byte(`
workloadOptions:
  llm:
    replicas: 2
    accelerator:
      perPod: split
      containers: [server]
      modes:
        - mode: nvidia
          supportMultiCards: true
          requiredGPUMemory: 48Gi
  web:
    replicas: 1
`)
	var got AppConfiguration
	if err := yaml.Unmarshal(src, &got); err != nil {
		t.Fatal(err)
	}
	if !got.WorkloadOptions.HasAccelerator() {
		t.Fatal("HasAccelerator = false, want true")
	}
	acc := got.WorkloadOptions["llm"].Accelerator
	if acc == nil || acc.PerPod != PerPodSplit || len(acc.Containers) != 1 || acc.Containers[0] != "server" {
		t.Fatalf("accelerator = %#v", acc)
	}
	if len(acc.Modes) != 1 || acc.Modes[0].Mode != "nvidia" || !acc.Modes[0].SupportMultiCards || acc.Modes[0].RequiredGPU != "48Gi" {
		t.Fatalf("modes = %#v", acc.Modes)
	}
	if got.WorkloadOptions["web"].Accelerator != nil {
		t.Fatal("web accelerator should be nil")
	}
	if (WorkloadOptions{"web": {}}).HasAccelerator() {
		t.Fatal("HasAccelerator = true without accelerator")
	}
}
