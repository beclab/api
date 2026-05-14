package v1alpha1

import (
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DefaultThirdLevelDomainConfig struct {
	AppName          string `json:"appName"`
	EntranceName     string `json:"entranceName"`
	ThirdLevelDomain string `json:"thirdLevelDomain"`
}

const (
	// AppApiVersionLabel marks an Application / ApplicationManager as a v3
	// (cluster-wide, single-CR) install. It is stamped at install time by
	// the v3 install handler and propagated by the Application controller.
	AppApiVersionLabel = "app.bytetrade.io/api-version"
	// AppVersionV3 is the value of AppApiVersionLabel for v3 apps.
	AppVersionV3 = "v3"

	// userSettingsKeyAuthLevel is the key inside Spec.UserSettings[user]
	// that stores a JSON blob mapping entrance name → auth level.
	userSettingsKeyAuthLevel = "authLevel"
)

// IsV3 reports whether the given object (Application or ApplicationManager)
// carries the v3 marker label.
func IsV3(o metav1.Object) bool {
	if o == nil {
		return false
	}
	return o.GetLabels()[AppApiVersionLabel] == AppVersionV3
}

// EffectiveSettings returns Spec.Settings overlaid with UserSettings[user]
// for v3 apps. For v1/v2 (no v3 label) it returns a copy of Spec.Settings
// as-is. The overlay is value-level on the top-level keys (e.g.
// "policy" / "customDomain" / "authLevel"); each value is a JSON blob
// keyed by entrance name that callers parse as today. Missing keys fall
// back to Spec.Settings. The returned map is always a fresh copy and
// never aliases the CR. Safe to call on a nil receiver.
func (app *Application) EffectiveSettings(user string) map[string]string {
	if app == nil {
		return map[string]string{}
	}
	out := make(map[string]string, len(app.Spec.Settings))
	for k, v := range app.Spec.Settings {
		out[k] = v
	}
	if !IsV3(app) || user == "" {
		return out
	}

	overlay, ok := app.Spec.UserSettings[user]
	if !ok {
		return out
	}
	for k, v := range overlay {
		out[k] = v
	}
	return out
}

// EffectiveEntrances returns a copy of Spec.Entrances with each entry's
// AuthLevel replaced by UserSettings[user]["authLevel"][name] when present.
// For v1/v2 apps it returns a copy of Spec.Entrances. The original CR is
// never mutated. A malformed authLevel overlay is ignored — callers fall
// back to the global Entrances rather than crashing because one user wrote
// junk into the CR. Safe to call on a nil receiver.
func (app *Application) EffectiveEntrances(user string) []Entrance {
	if app == nil {
		return nil
	}
	out := make([]Entrance, len(app.Spec.Entrances))
	copy(out, app.Spec.Entrances)
	if !IsV3(app) || user == "" {
		return out
	}
	overlay, ok := app.Spec.UserSettings[user]
	if !ok {
		return out
	}
	authBlob, ok := overlay[userSettingsKeyAuthLevel]
	if !ok || authBlob == "" {
		return out
	}
	var perEntrance map[string]string
	if err := json.Unmarshal([]byte(authBlob), &perEntrance); err != nil {
		return out
	}
	for i := range out {
		if lvl, ok := perEntrance[out[i].Name]; ok && lvl != "" {
			out[i].AuthLevel = lvl
		}
	}
	return out
}
