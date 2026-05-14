package v1alpha1

import (
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// helper to build an Application with the v3 label.
func newV3App(spec ApplicationSpec) *Application {
	return &Application{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "demo",
			Labels: map[string]string{AppApiVersionLabel: AppVersionV3},
		},
		Spec: spec,
	}
}

// helper to build an Application without the v3 label (v1/v2).
func newV1App(spec ApplicationSpec) *Application {
	return &Application{
		ObjectMeta: metav1.ObjectMeta{Name: "demo"},
		Spec:       spec,
	}
}

func TestIsV3(t *testing.T) {
	tests := []struct {
		name string
		obj  metav1.Object
		want bool
	}{
		{
			name: "nil object",
			obj:  nil,
			want: false,
		},
		{
			name: "no labels",
			obj:  &Application{ObjectMeta: metav1.ObjectMeta{Name: "x"}},
			want: false,
		},
		{
			name: "wrong label value",
			obj: &Application{ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{AppApiVersionLabel: "v2"},
			}},
			want: false,
		},
		{
			name: "v3 label on Application",
			obj: &Application{ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{AppApiVersionLabel: AppVersionV3},
			}},
			want: true,
		},
		{
			name: "v3 label on ApplicationManager",
			obj: &ApplicationManager{ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{AppApiVersionLabel: AppVersionV3},
			}},
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsV3(tc.obj); got != tc.want {
				t.Fatalf("IsV3() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestEffectiveSettings(t *testing.T) {
	t.Run("nil app returns empty non-nil map", func(t *testing.T) {
		got := (*Application)(nil).EffectiveSettings("alice")
		if got == nil || len(got) != 0 {
			t.Fatalf("EffectiveSettings(nil) = %v, want empty map", got)
		}
	})

	t.Run("v1/v2 app returns Settings as-is and is a copy", func(t *testing.T) {
		app := newV1App(ApplicationSpec{
			Settings: map[string]string{
				"customDomain": `{"e1":{"third_party_domain":"a.example.com"}}`,
				"policy":       `{"e1":{"default_policy":"public"}}`,
			},
			UserSettings: map[string]map[string]string{
				"alice": {"customDomain": "should-be-ignored"},
			},
		})
		got := app.EffectiveSettings("alice")
		want := map[string]string{
			"customDomain": `{"e1":{"third_party_domain":"a.example.com"}}`,
			"policy":       `{"e1":{"default_policy":"public"}}`,
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("EffectiveSettings v1/v2 = %v, want %v", got, want)
		}

		// Mutating the returned map must not affect the CR.
		got["customDomain"] = "mutated"
		if app.Spec.Settings["customDomain"] == "mutated" {
			t.Fatal("EffectiveSettings returned a map aliasing Spec.Settings")
		}
	})

	t.Run("v3 app with empty user returns global Settings", func(t *testing.T) {
		app := newV3App(ApplicationSpec{
			Settings: map[string]string{"title": "Demo"},
			UserSettings: map[string]map[string]string{
				"alice": {"customDomain": `{"e1":{"third_party_domain":"a.example.com"}}`},
			},
		})
		got := app.EffectiveSettings("")
		want := map[string]string{"title": "Demo"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("EffectiveSettings v3 empty user = %v, want %v", got, want)
		}
	})

	t.Run("v3 app with no overlay returns global Settings", func(t *testing.T) {
		app := newV3App(ApplicationSpec{
			Settings: map[string]string{"title": "Demo"},
			UserSettings: map[string]map[string]string{
				"alice": {"customDomain": "a-only"},
			},
		})
		got := app.EffectiveSettings("bob")
		want := map[string]string{"title": "Demo"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("EffectiveSettings v3 no-overlay = %v, want %v", got, want)
		}
	})

	t.Run("v3 app overlays user-specific keys on top of Settings", func(t *testing.T) {
		app := newV3App(ApplicationSpec{
			Settings: map[string]string{
				"title":        "Demo",
				"customDomain": `{"e1":{"third_party_domain":"global.example.com"}}`,
			},
			UserSettings: map[string]map[string]string{
				"alice": {
					"customDomain": `{"e1":{"third_party_domain":"alice.example.com"}}`,
					"policy":       `{"e1":{"default_policy":"private"}}`,
				},
				"bob": {
					"customDomain": `{"e1":{"third_party_domain":"bob.example.com"}}`,
				},
			},
		})

		gotAlice := app.EffectiveSettings("alice")
		wantAlice := map[string]string{
			"title":        "Demo",
			"customDomain": `{"e1":{"third_party_domain":"alice.example.com"}}`,
			"policy":       `{"e1":{"default_policy":"private"}}`,
		}
		if !reflect.DeepEqual(gotAlice, wantAlice) {
			t.Fatalf("EffectiveSettings(alice) = %v, want %v", gotAlice, wantAlice)
		}

		gotBob := app.EffectiveSettings("bob")
		wantBob := map[string]string{
			"title":        "Demo",
			"customDomain": `{"e1":{"third_party_domain":"bob.example.com"}}`,
		}
		if !reflect.DeepEqual(gotBob, wantBob) {
			t.Fatalf("EffectiveSettings(bob) = %v, want %v", gotBob, wantBob)
		}

		// Per-user views must not leak across users (verifies the
		// returned map is freshly allocated, not shared).
		gotAlice["customDomain"] = "tampered"
		if app.EffectiveSettings("bob")["customDomain"] == "tampered" {
			t.Fatal("EffectiveSettings results alias across users")
		}
		if app.Spec.Settings["customDomain"] != `{"e1":{"third_party_domain":"global.example.com"}}` {
			t.Fatal("EffectiveSettings mutated Spec.Settings")
		}
	})
}

func TestEffectiveEntrances(t *testing.T) {
	baseEntrances := []Entrance{
		{Name: "e1", Host: "h1", Port: 80, AuthLevel: "public"},
		{Name: "e2", Host: "h2", Port: 81, AuthLevel: "private"},
	}

	t.Run("nil app returns nil", func(t *testing.T) {
		if got := (*Application)(nil).EffectiveEntrances("alice"); got != nil {
			t.Fatalf("EffectiveEntrances(nil) = %v, want nil", got)
		}
	})

	t.Run("v1/v2 returns Entrances copy and ignores UserSettings", func(t *testing.T) {
		app := newV1App(ApplicationSpec{
			Entrances: baseEntrances,
			UserSettings: map[string]map[string]string{
				"alice": {"authLevel": `{"e1":"private"}`},
			},
		})
		got := app.EffectiveEntrances("alice")
		if !reflect.DeepEqual(got, baseEntrances) {
			t.Fatalf("EffectiveEntrances v1/v2 = %v, want %v", got, baseEntrances)
		}

		// Must not alias the CR slice.
		got[0].AuthLevel = "mutated"
		if app.Spec.Entrances[0].AuthLevel == "mutated" {
			t.Fatal("EffectiveEntrances returned a slice aliasing Spec.Entrances")
		}
	})

	t.Run("v3 app overlays AuthLevel from UserSettings[user][authLevel]", func(t *testing.T) {
		app := newV3App(ApplicationSpec{
			Entrances: baseEntrances,
			UserSettings: map[string]map[string]string{
				"alice": {"authLevel": `{"e1":"private","e2":"public"}`},
				"bob":   {"authLevel": `{"e1":"private"}`},
			},
		})

		gotAlice := app.EffectiveEntrances("alice")
		wantAlice := []Entrance{
			{Name: "e1", Host: "h1", Port: 80, AuthLevel: "private"},
			{Name: "e2", Host: "h2", Port: 81, AuthLevel: "public"},
		}
		if !reflect.DeepEqual(gotAlice, wantAlice) {
			t.Fatalf("EffectiveEntrances(alice) = %v, want %v", gotAlice, wantAlice)
		}

		// bob only overrode e1; e2 must fall back to global.
		gotBob := app.EffectiveEntrances("bob")
		wantBob := []Entrance{
			{Name: "e1", Host: "h1", Port: 80, AuthLevel: "private"},
			{Name: "e2", Host: "h2", Port: 81, AuthLevel: "private"},
		}
		if !reflect.DeepEqual(gotBob, wantBob) {
			t.Fatalf("EffectiveEntrances(bob) = %v, want %v", gotBob, wantBob)
		}

		// Global Entrances must be untouched after either call.
		if app.Spec.Entrances[0].AuthLevel != "public" || app.Spec.Entrances[1].AuthLevel != "private" {
			t.Fatalf("Spec.Entrances mutated: %+v", app.Spec.Entrances)
		}
	})

	t.Run("v3 app empty user returns global Entrances", func(t *testing.T) {
		app := newV3App(ApplicationSpec{
			Entrances: baseEntrances,
			UserSettings: map[string]map[string]string{
				"alice": {"authLevel": `{"e1":"private"}`},
			},
		})
		got := app.EffectiveEntrances("")
		if !reflect.DeepEqual(got, baseEntrances) {
			t.Fatalf("EffectiveEntrances v3 empty user = %v, want %v", got, baseEntrances)
		}
	})

	t.Run("v3 app missing authLevel overlay returns global Entrances", func(t *testing.T) {
		app := newV3App(ApplicationSpec{
			Entrances: baseEntrances,
			UserSettings: map[string]map[string]string{
				"alice": {"customDomain": "x"},
			},
		})
		got := app.EffectiveEntrances("alice")
		if !reflect.DeepEqual(got, baseEntrances) {
			t.Fatalf("EffectiveEntrances missing authLevel = %v, want %v", got, baseEntrances)
		}
	})

	t.Run("v3 app malformed authLevel JSON falls back to globals", func(t *testing.T) {
		app := newV3App(ApplicationSpec{
			Entrances: baseEntrances,
			UserSettings: map[string]map[string]string{
				"alice": {"authLevel": `{not valid json`},
			},
		})
		got := app.EffectiveEntrances("alice")
		if !reflect.DeepEqual(got, baseEntrances) {
			t.Fatalf("EffectiveEntrances malformed = %v, want %v", got, baseEntrances)
		}
	})

	t.Run("v3 app empty string AuthLevel value is ignored", func(t *testing.T) {
		app := newV3App(ApplicationSpec{
			Entrances: baseEntrances,
			UserSettings: map[string]map[string]string{
				"alice": {"authLevel": `{"e1":""}`},
			},
		})
		got := app.EffectiveEntrances("alice")
		if got[0].AuthLevel != "public" {
			t.Fatalf("EffectiveEntrances empty value override = %q, want %q", got[0].AuthLevel, "public")
		}
	})

	t.Run("v3 app unknown entrance name in overlay is silently skipped", func(t *testing.T) {
		app := newV3App(ApplicationSpec{
			Entrances: baseEntrances,
			UserSettings: map[string]map[string]string{
				"alice": {"authLevel": `{"e1":"private","does-not-exist":"public"}`},
			},
		})
		got := app.EffectiveEntrances("alice")
		want := []Entrance{
			{Name: "e1", Host: "h1", Port: 80, AuthLevel: "private"},
			{Name: "e2", Host: "h2", Port: 81, AuthLevel: "private"},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("EffectiveEntrances unknown entrance = %v, want %v", got, want)
		}
	})
}
