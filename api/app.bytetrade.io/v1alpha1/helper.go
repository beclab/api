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
	ThirdLevelDomain string `json:"thirdLevelDomain"`
}

const (
	// AppApiVersionLabel marks an Application / ApplicationManager as a v3
	// (cluster-wide, single-CR) install. It is stamped at install time by
	// the v3 install handler and propagated by the Application controller.
	AppApiVersionLabel = "app.bytetrade.io/api-version"
	// AppVersionV3 is the value of AppApiVersionLabel for v3 apps.
	AppVersionV3 = "v3"

	AppVersionV1 = "v1"

	// userSettingsKeyAuthLevel is the override key that stores a JSON blob
	// mapping entrance name → auth level. It lives in Spec.UserSettings[user]
	// for shared apps and in Spec.Settings for non-shared apps.
	userSettingsKeyAuthLevel = "authLevel"

	// userSettingsKeyEntranceOverrides is the override key that stores a JSON
	// blob mapping entrance name → EntranceOverride (the user-editable non-auth
	// fields of an entrance). Same location rules as userSettingsKeyAuthLevel.
	userSettingsKeyEntranceOverrides = "entranceOverrides"

	// userSettingsKeyAddedEntrances is the override key that stores a JSON array
	// ([]Entrance) of entrances added on top of the chart-derived
	// Spec.Entrances. Same location rules as userSettingsKeyAuthLevel.
	userSettingsKeyAddedEntrances = "addedEntrances"

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

// EntranceOverride carries the user-editable, non-auth fields of an Entrance.
// All fields are pointers so an unset field leaves the base (chart) value
// untouched — this distinguishes "not overridden" from a zero value, which
// matters for the bool fields where false is a legitimate override. AuthLevel
// is intentionally excluded: it keeps its own dedicated UserSettings key
// (userSettingsKeyAuthLevel) for backward compatibility.
type EntranceOverride struct {
	Title           *string `json:"title,omitempty"`
	Icon            *string `json:"icon,omitempty"`
	Invisible       *bool   `json:"invisible,omitempty"`
	OpenMethod      *string `json:"openMethod,omitempty"`
	WindowPushState *bool   `json:"windowPushState,omitempty"`
	URL             *string `json:"url,omitempty"`
}

// applyTo overwrites the set fields of the override onto the entrance.
func (o EntranceOverride) applyTo(e *Entrance) {
	if o.Title != nil {
		e.Title = *o.Title
	}
	if o.Icon != nil {
		e.Icon = *o.Icon
	}
	if o.Invisible != nil {
		e.Invisible = *o.Invisible
	}
	if o.OpenMethod != nil {
		e.OpenMethod = *o.OpenMethod
	}
	if o.WindowPushState != nil {
		e.WindowPushState = *o.WindowPushState
	}
	if o.URL != nil {
		e.URL = *o.URL
	}
}

// overrideView returns the map that holds the user/background overrides for the
// given caller, and whether any override view applies. Shared (v3) apps keep a
// per-user override map at Spec.UserSettings[user]; a caller-less lookup ("")
// has no view. Non-shared (v1/v2/v3) apps keep their overrides in the app-global
// Spec.Settings itself, so that map is returned regardless of the caller. Safe
// on nil.
func (app *Application) overrideView(user string) (map[string]string, bool) {
	if app == nil {
		return nil, false
	}
	if IsShared(app) {
		if user == "" {
			return nil, false
		}
		overlay, ok := app.Spec.UserSettings[user]
		return overlay, ok
	}
	return app.Spec.Settings, true
}

// EffectiveSettings returns Spec.Settings, overlaid with the per-user
// Spec.UserSettings[user] entry for shared (v3) apps. Non-shared apps store
// their overrides directly in Spec.Settings, so their effective settings are
// just a copy of Spec.Settings. The returned map is always a fresh copy and
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

// EffectiveEntrances returns a copy of Spec.Entrances (the chart-derived base)
// with the applicable overrides applied: user-added entrances are appended,
// per-entrance AuthLevel is replaced from the "authLevel" blob, and remaining
// user-editable fields are replaced from the "entranceOverrides" blob. The
// override source is the per-user Spec.UserSettings[user] map for shared (v3)
// apps and Spec.Settings for non-shared apps. The original CR is never mutated.
// Any malformed override blob is ignored — callers fall back to the base rather
// than crashing because someone wrote junk into the CR. Safe to call on a nil
// receiver.
func (app *Application) EffectiveEntrances(user string) []Entrance {
	if app == nil {
		return nil
	}
	out := make([]Entrance, len(app.Spec.Entrances))
	copy(out, app.Spec.Entrances)

	overlay, ok := app.overrideView(user)
	if !ok || overlay == nil {
		return out
	}

	// Append user-added entrances first so the per-name overlays below can
	// also apply to them (e.g. an authLevel change on an added entrance).
	if added := overlay[userSettingsKeyAddedEntrances]; added != "" {
		var addedEntrances []Entrance
		if err := json.Unmarshal([]byte(added), &addedEntrances); err == nil {
			out = append(out, addedEntrances...)
		}
	}

	if authBlob := overlay[userSettingsKeyAuthLevel]; authBlob != "" {
		var perEntrance map[string]string
		if err := json.Unmarshal([]byte(authBlob), &perEntrance); err == nil {
			for i := range out {
				if lvl, ok := perEntrance[out[i].Name]; ok && lvl != "" {
					out[i].AuthLevel = lvl
				}
			}
		}
	}

	if ovBlob := overlay[userSettingsKeyEntranceOverrides]; ovBlob != "" {
		var perEntrance map[string]EntranceOverride
		if err := json.Unmarshal([]byte(ovBlob), &perEntrance); err == nil {
			for i := range out {
				if ov, ok := perEntrance[out[i].Name]; ok {
					ov.applyTo(&out[i])
				}
			}
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

// SharedEntranceIDV2 returns the id of a single shared entrance, honouring the
// single-entrance rule. It mirrors EntranceID but the prefix is
// md5(appid + "shared")[:8] instead of the bare appid.
func SharedEntranceIDV2(appid string, entranceIndex, entranceCount int) string {
	prefix := SharedEntrancePrefix(appid)
	return fmt.Sprintf("%s%d", prefix, entranceIndex)
}

// SharedEntranceIDV2 returns the id of this shared entrance for the given appid,
// honouring the single-entrance rule. See SharedEntranceID for the id format.
func (e Entrance) SharedEntranceIDV2(appid string, entranceIndex, entranceCount int) string {
	return SharedEntranceIDV2(appid, entranceIndex, entranceCount)
}

// SharedForZone returns a copy of this shared entrance with its URL rewritten
// to "<sharedEntranceID>.<zone>" for the given appid. The receiver is never
// mutated.
func (e Entrance) SharedForZone(appid, zone string, entranceIndex, entranceCount int) Entrance {
	out := e
	out.URL = fmt.Sprintf("%s.%s", e.SharedEntranceID(appid, entranceIndex, entranceCount), zone)
	return out
}

// SharedForZoneV2 returns a copy of this shared entrance with its URL rewritten
// to "<sharedEntranceID>.<zone>" for the given appid. The receiver is never
// mutated.
func (e Entrance) SharedForZoneV2(appid, zone string, entranceIndex, entranceCount int) Entrance {
	out := e
	out.URL = fmt.Sprintf("%s.%s", e.SharedEntranceIDV2(appid, entranceIndex, entranceCount), zone)
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

// SharedForZoneV2 returns a copy of the list with each entry's URL rewritten to
// "<sharedEntranceID>.<zone>" for the given appid. The receiver is never
// mutated.
func (es Entrances) SharedForZoneV2(appid, zone string) Entrances {
	out := make(Entrances, len(es))
	for i := range es {
		out[i] = es[i].SharedForZoneV2(appid, zone, i, len(es))
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

// ThirdLevelCusDomainURLs returns the configured third-level domain
// prefixes for every entrance of the application, each suffixed with
// ".<zone>" to form a full host. When zone is empty the bare prefixes are
// returned. The lookup uses Spec.Settings overlaid with the install owner's
// UserSettings, so shared v3 apps see the owner's per-entrance overrides.
// Safe to call on a nil receiver.
func (app *Application) ThirdLevelCusDomainURLs(zone string, owner string) []string {
	if app == nil {
		return nil
	}
	effectiveEntrances := app.EffectiveEntrances(owner)
	if len(effectiveEntrances) == 0 {
		return nil
	}
	customDomainEntrancesMap := settingsEntranceMap(app.EffectiveSettings(owner), settingsKeyCustomDomain)

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
