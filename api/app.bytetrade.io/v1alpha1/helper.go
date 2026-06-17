package v1alpha1

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DefaultThirdLevelDomainConfig struct {
	AppName          string `json:"appName"`
	EntranceName     string `json:"entranceName"`
	ThirdLevelDomain string `json:"third_level_domain"`
}

const (
	// AppApiVersionLabel marks an Application / ApplicationManager as a v3
	// (cluster-wide, single-CR) install. It is stamped at install time by
	// the v3 install handler and propagated by the Application controller.
	AppApiVersionLabel = "app.bytetrade.io/api-version"
	// AppVersionV3 is the value of AppApiVersionLabel for v3 apps.
	AppVersionV3 = "v3"

	AppVersionV1 = "v1"

	// userSettingsKeyAuthLevel is the key inside Spec.UserSettings[user]
	// that stores a JSON blob mapping entrance name → auth level.
	userSettingsKeyAuthLevel = "authLevel"

	settingsKeyCustomDomain = "customDomain"
	// settingsKeyDefaultThirdLevelDomainConfig stores per-app default third-level
	// domain overrides as a JSON array of DefaultThirdLevelDomainConfig.
	settingsKeyDefaultThirdLevelDomainConfig = "defaultThirdLevelDomainConfig"
	// settingsCustomDomainThirdLevelDomain is the per-entrance key inside the
	// customDomain JSON blob for a user-defined third-level domain prefix.
	settingsCustomDomainThirdLevelDomain = "third_level_domain"

	ApplicationAuthLevelPublic = "public"

	AppSharedLabel = "app.bytetrade.io/app-shared"
	AppSharedTrue  = "true"
)

// IsV3 reports whether the given object (Application or ApplicationManager)
// carries the v3 marker label.
func IsV3(o metav1.Object) bool {
	if o == nil {
		return false
	}
	return o.GetLabels()[AppApiVersionLabel] == AppVersionV3
}

func IsShared(o metav1.Object) bool {
	if o == nil {
		return false
	}
	return o.GetLabels()[AppSharedLabel] == AppSharedTrue
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
	if !IsShared(app) || user == "" {
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
	if !IsShared(app) || user == "" {
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

// EntranceID builds the identifier of a single entrance from the app id and
// the entrance's position. It is the single source of truth for the shared
// rule: when an application exposes only one entrance the bare appid is
// returned, otherwise the entrance index is appended (e.g. appid "abc123" at
// index 2 -> "abc1232"). entranceIndex is 0-based and entranceCount is the
// total number of entrances of the application.
func EntranceID(appid string, entranceIndex, entranceCount int) string {
	if entranceCount <= 1 {
		return appid
	}
	return fmt.Sprintf("%s%d", appid, entranceIndex)
}

// EntranceID returns the id of this entrance for the given appid, honouring
// the single-entrance rule. entranceIndex is the 0-based position of the
// entrance within its application and entranceCount is the total number of
// entrances. Unlike the application-level helpers, a bare Entrance does not
// carry the appid, so it must be supplied by the caller.
func (e Entrance) EntranceID(appid string, entranceIndex, entranceCount int) string {
	return EntranceID(appid, entranceIndex, entranceCount)
}

// ForZone returns a copy of this entrance with its URL rewritten to
// "<entranceID>.<zone>" for the given appid, honouring the single-entrance
// rule. The receiver is never mutated.
func (e Entrance) ForZone(appid, zone string, entranceIndex, entranceCount int) Entrance {
	out := e
	out.URL = fmt.Sprintf("%s.%s", e.EntranceID(appid, entranceIndex, entranceCount), zone)
	return out
}

// Entrances is a named type over []Entrance that provides bulk helpers
// honouring the single-entrance rule. Callers holding a plain slice can
// convert with Entrances(s). The slice length supplies the entrance count, so
// the single-entrance rule is applied automatically.
type Entrances []Entrance

// EntranceIDs returns the id of every entrance in the list for the given
// appid, preserving order and honouring the single-entrance rule.
func (es Entrances) EntranceIDs(appid string) []string {
	ids := make([]string, len(es))
	for i := range es {
		ids[i] = es[i].EntranceID(appid, i, len(es))
	}
	return ids
}

// ForZone returns a copy of the list with each entry's URL rewritten to
// "<entranceID>.<zone>" for the given appid, honouring the single-entrance
// rule. The receiver is never mutated.
func (es Entrances) ForZone(appid, zone string) Entrances {
	out := make(Entrances, len(es))
	for i := range es {
		out[i] = es[i].ForZone(appid, zone, i, len(es))
	}
	return out
}

// EntranceIDs returns the entrance id of every entrance of the application,
// preserving entrance order and honouring the single-entrance rule. Safe to
// call on a nil receiver.
func (app *Application) EntranceIDs() []string {
	if app == nil {
		return nil
	}
	return Entrances(app.Spec.Entrances).EntranceIDs(app.Spec.Appid)
}

// EntrancesForZone returns a copy of the application entrances with each URL
// rewritten to "<entranceID>.<zone>", where entranceID honours the
// single-entrance rule. The original CR is never mutated. Safe to call on a
// nil receiver.
func (app *Application) EntrancesForZone(zone string) []Entrance {
	if app == nil {
		return nil
	}
	return Entrances(app.Spec.Entrances).ForZone(app.Spec.Appid, zone)
}

// SharedEntrancePrefix returns the 8-char prefix used to build shared entrance
// ids: the first 8 hex chars of md5(appid + "shared"). This is the only
// difference from the regular entrance id, which uses the bare appid.
func SharedEntrancePrefix(appid string) string {
	hash := md5.Sum([]byte(appid + "shared"))
	return hex.EncodeToString(hash[:])[:8]
}

// SharedEntranceID returns the id of a single shared entrance, honouring the
// single-entrance rule. It mirrors EntranceID but the prefix is
// md5(appid + "shared")[:8] instead of the bare appid.
func SharedEntranceID(appid string, entranceIndex, entranceCount int) string {
	prefix := SharedEntrancePrefix(appid)
	if entranceCount <= 1 {
		return prefix
	}
	return fmt.Sprintf("%s%d", prefix, entranceIndex)
}

// SharedEntranceID returns the id of this shared entrance for the given appid,
// honouring the single-entrance rule. See SharedEntranceID for the id format.
func (e Entrance) SharedEntranceID(appid string, entranceIndex, entranceCount int) string {
	return SharedEntranceID(appid, entranceIndex, entranceCount)
}

// SharedForZone returns a copy of this shared entrance with its URL rewritten
// to "<sharedEntranceID>.<zone>" for the given appid. The receiver is never
// mutated.
func (e Entrance) SharedForZone(appid, zone string, entranceIndex, entranceCount int) Entrance {
	out := e
	out.URL = fmt.Sprintf("%s.%s", e.SharedEntranceID(appid, entranceIndex, entranceCount), zone)
	return out
}

// SharedEntranceIDs returns the shared entrance id of every entrance in the
// list for the given appid, preserving order and honouring the single-entrance
// rule.
func (es Entrances) SharedEntranceIDs(appid string) []string {
	ids := make([]string, len(es))
	for i := range es {
		ids[i] = es[i].SharedEntranceID(appid, i, len(es))
	}
	return ids
}

// SharedForZone returns a copy of the list with each entry's URL rewritten to
// "<sharedEntranceID>.<zone>" for the given appid. The receiver is never
// mutated.
func (es Entrances) SharedForZone(appid, zone string) Entrances {
	out := make(Entrances, len(es))
	for i := range es {
		out[i] = es[i].SharedForZone(appid, zone, i, len(es))
	}
	return out
}

// SharedEntranceIDs returns the shared entrance id of every entry in
// Spec.SharedEntrances, honouring the single-entrance rule. Safe to call on a
// nil receiver.
func (app *Application) SharedEntranceIDs() []string {
	if app == nil {
		return nil
	}
	return Entrances(app.Spec.SharedEntrances).SharedEntranceIDs(app.Spec.Appid)
}

// SharedEntrancesForZone returns a copy of Spec.SharedEntrances with each URL
// rewritten to "<sharedEntranceID>.<zone>". The original CR is never mutated.
// Safe to call on a nil receiver.
func (app *Application) SharedEntrancesForZone(zone string) []Entrance {
	if app == nil {
		return nil
	}
	return Entrances(app.Spec.SharedEntrances).SharedForZone(app.Spec.Appid, zone)
}

// settingsEntranceMap parses settings[key] as a JSON object mapping entrance
// name → key/value pairs. Malformed or missing JSON yields nil.
func settingsEntranceMap(settings map[string]string, key string) map[string]map[string]string {
	blob, ok := settings[key]
	if !ok || blob == "" {
		return nil
	}
	var out map[string]map[string]string
	if err := json.Unmarshal([]byte(blob), &out); err != nil {
		return nil
	}
	return out
}

// ThirdLevelCusDomainPrefixes returns the configured third-level domain
// prefixes for every entrance of the application, each suffixed with
// ".<zone>" to form a full host. When zone is empty the bare prefixes are
// returned. The lookup uses Spec.Settings overlaid with the install owner's
// UserSettings, so shared v3 apps see the owner's per-entrance overrides.
// Safe to call on a nil receiver.
func (app *Application) ThirdLevelCusDomainPrefixes(zone string) []string {
	if app == nil {
		return nil
	}
	effectiveEntrances := app.EffectiveEntrances(app.Spec.Owner)
	if len(effectiveEntrances) == 0 {
		return nil
	}
	customDomainEntrancesMap := settingsEntranceMap(app.EffectiveSettings(app.Spec.Owner), settingsKeyCustomDomain)

	var out []string
	for _, entrance := range effectiveEntrances {
		cdEntrance, ok := customDomainEntrancesMap[entrance.Name]
		if !ok {
			continue
		}
		entrancePrefix := cdEntrance[settingsCustomDomainThirdLevelDomain]
		if entrancePrefix == "" {
			continue
		}
		if zone == "" {
			continue
		}
		out = append(out, fmt.Sprintf("%s.%s", entrancePrefix, zone))
	}
	return out
}

// EntrancesWithZone returns a copy of Spec.Entrances with each URL rewritten
// for the given zone. When zone is empty the entrances are returned unchanged.
// defaultThirdLevelDomainConfig in Spec.Settings can override individual
// entrance URLs. The original CR is never mutated. Safe to call on a nil
// receiver.
func (app *Application) EntrancesWithZone(zone string) ([]Entrance, error) {
	if app == nil {
		return nil, nil
	}
	out := make([]Entrance, len(app.Spec.Entrances))
	copy(out, app.Spec.Entrances)
	if zone == "" {
		return out, nil
	}

	var appDomainConfigs []DefaultThirdLevelDomainConfig
	if defaultThirdLevelDomainConfig, ok := app.Spec.Settings[settingsKeyDefaultThirdLevelDomainConfig]; ok && defaultThirdLevelDomainConfig != "" {
		if err := json.Unmarshal([]byte(defaultThirdLevelDomainConfig), &appDomainConfigs); err != nil {
			return nil, err
		}
	}

	appid := strings.ToLower(strings.TrimSpace(app.Spec.Appid))
	n := len(out)
	if n == 1 {
		out[0] = out[0].ForZone(appid, zone, 0, 1)
		return out, nil
	}

	entrancesForZone := Entrances(out).ForZone(appid, zone)
	for i := range entrancesForZone {
		out[i] = entrancesForZone[i]
		for _, adc := range appDomainConfigs {
			if adc.AppName == app.Spec.Name && adc.EntranceName == out[i].Name && adc.ThirdLevelDomain != "" {
				out[i].URL = fmt.Sprintf("%s.%s", adc.ThirdLevelDomain, zone)
			}
		}
	}
	return out, nil
}

// GenEntranceURLs returns a copy of Spec.Entrances with URLs filled from the
// install owner's zone annotation (bytetrade.io/zone). Zone lookup errors are
// ignored and the entrances are returned unchanged, matching legacy Provider
// behaviour. Malformed defaultThirdLevelDomainConfig returns an error. The
// original CR is never mutated. Safe to call on a nil receiver.
func (app *Application) GenEntranceURLs(zone string) ([]Entrance, error) {
	if app == nil {
		return nil, nil
	}
	return app.EntrancesWithZone(zone)
}
