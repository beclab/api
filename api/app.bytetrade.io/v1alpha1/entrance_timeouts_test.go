package v1alpha1

import (
	"encoding/json"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func TestEntranceTimeoutsDeepCopy(t *testing.T) {
	idle := metav1.Duration{Duration: 15 * time.Minute}
	total := metav1.Duration{Duration: time.Hour}
	in := Entrance{
		Name: "api",
		Host: "api",
		Port: 8080,
		Timeouts: &EntranceTimeouts{
			Idle:  &idle,
			Total: &total,
		},
		MaxBodySize: "100Mi",
	}

	out := in.DeepCopy()
	if out.Timeouts == nil || out.Timeouts.Idle == nil || out.Timeouts.Total == nil {
		t.Fatalf("DeepCopy dropped timeouts: %+v", out.Timeouts)
	}
	if out.Timeouts == in.Timeouts || out.Timeouts.Idle == in.Timeouts.Idle || out.Timeouts.Total == in.Timeouts.Total {
		t.Fatal("DeepCopy shared timeout pointers")
	}
	out.Timeouts.Idle.Duration = time.Second
	out.MaxBodySize = "1Mi"
	if in.Timeouts.Idle.Duration != 15*time.Minute {
		t.Fatalf("mutating copy changed original idle: %s", in.Timeouts.Idle.Duration)
	}
	if in.MaxBodySize != "100Mi" {
		t.Fatalf("mutating copy changed original maxBodySize: %q", in.MaxBodySize)
	}
}

func TestEntranceTimeoutsYAMLAndJSONRoundTrip(t *testing.T) {
	const raw = `
name: api
host: api
port: 8080
timeouts:
  idle: 900s
  total: 3600s
maxBodySize: 100Mi
`
	var fromYAML Entrance
	if err := yaml.Unmarshal([]byte(raw), &fromYAML); err != nil {
		t.Fatalf("yaml.Unmarshal: %v", err)
	}
	if fromYAML.Timeouts == nil || fromYAML.Timeouts.Idle == nil || fromYAML.Timeouts.Total == nil {
		t.Fatalf("yaml timeouts = %+v", fromYAML.Timeouts)
	}
	if fromYAML.Timeouts.Idle.Duration != 900*time.Second || fromYAML.Timeouts.Total.Duration != 3600*time.Second {
		t.Fatalf("yaml durations idle=%s total=%s", fromYAML.Timeouts.Idle.Duration, fromYAML.Timeouts.Total.Duration)
	}
	if fromYAML.MaxBodySize != "100Mi" {
		t.Fatalf("yaml maxBodySize = %q", fromYAML.MaxBodySize)
	}

	b, err := json.Marshal(fromYAML)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	var fromJSON Entrance
	if err := json.Unmarshal(b, &fromJSON); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if fromJSON.Timeouts == nil || fromJSON.Timeouts.Total == nil || fromJSON.Timeouts.Total.Duration != 3600*time.Second {
		t.Fatalf("json timeouts = %+v", fromJSON.Timeouts)
	}
	if fromJSON.MaxBodySize != "100Mi" {
		t.Fatalf("json maxBodySize = %q", fromJSON.MaxBodySize)
	}
}

func TestEffectiveEntrancesDoesNotAliasTimeouts(t *testing.T) {
	idle := metav1.Duration{Duration: 10 * time.Minute}
	app := &Application{
		Spec: ApplicationSpec{
			Entrances: []Entrance{{
				Name: "web",
				Host: "web",
				Port: 80,
				Timeouts: &EntranceTimeouts{Idle: &idle},
			}},
		},
	}
	got := app.EffectiveEntrances("")
	if len(got) != 1 || got[0].Timeouts == nil || got[0].Timeouts.Idle == nil {
		t.Fatalf("EffectiveEntrances = %+v", got)
	}
	if got[0].Timeouts == app.Spec.Entrances[0].Timeouts || got[0].Timeouts.Idle == app.Spec.Entrances[0].Timeouts.Idle {
		t.Fatal("EffectiveEntrances shared Timeouts with Spec")
	}
	got[0].Timeouts.Idle.Duration = time.Second
	if app.Spec.Entrances[0].Timeouts.Idle.Duration != 10*time.Minute {
		t.Fatalf("Spec idle mutated to %s", app.Spec.Entrances[0].Timeouts.Idle.Duration)
	}
}

func TestEntrancesWithZoneDoesNotAliasTimeouts(t *testing.T) {
	total := metav1.Duration{Duration: 30 * time.Minute}
	app := &Application{
		Spec: ApplicationSpec{
			Appid: "demo1234",
			Entrances: []Entrance{{
				Name: "web",
				Host: "web",
				Port: 80,
				Timeouts: &EntranceTimeouts{Total: &total},
			}},
		},
	}
	got, err := app.EntrancesWithZone("olares.com")
	if err != nil {
		t.Fatalf("EntrancesWithZone: %v", err)
	}
	if len(got) != 1 || got[0].Timeouts == nil || got[0].Timeouts.Total == nil {
		t.Fatalf("EntrancesWithZone = %+v", got)
	}
	got[0].Timeouts.Total.Duration = time.Second
	if app.Spec.Entrances[0].Timeouts.Total.Duration != 30*time.Minute {
		t.Fatalf("Spec total mutated to %s", app.Spec.Entrances[0].Timeouts.Total.Duration)
	}
}

func TestForZoneDoesNotAliasTimeouts(t *testing.T) {
	idle := metav1.Duration{Duration: 5 * time.Minute}
	in := Entrance{
		Name:     "web",
		Host:     "web",
		Port:     80,
		Timeouts: &EntranceTimeouts{Idle: &idle},
	}
	out := in.ForZone("demo1234", "olares.com", 0, 1)
	if out.Timeouts == nil || out.Timeouts.Idle == nil {
		t.Fatalf("ForZone dropped timeouts: %+v", out.Timeouts)
	}
	out.Timeouts.Idle.Duration = time.Second
	if in.Timeouts.Idle.Duration != 5*time.Minute {
		t.Fatalf("ForZone aliased idle: %s", in.Timeouts.Idle.Duration)
	}
}
