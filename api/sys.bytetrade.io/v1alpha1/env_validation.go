package v1alpha1

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/mail"
	"net/url"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/validation"
)

func (e *EnvVarSpec) ValidateValue(value string) error {
	if value == "" {
		return nil
	}

	// Memoize the remote options fetch so a multi-select value with N items
	// triggers at most one HTTP request instead of one per item.
	var (
		remoteFetched bool
		remoteOptions []EnvValueOptionItem
		remoteErr     error
	)
	getRemote := func() ([]EnvValueOptionItem, error) {
		if !remoteFetched {
			remoteOptions, remoteErr = fetchRemoteOptions(e.RemoteOptions)
			remoteFetched = true
		}
		return remoteOptions, remoteErr
	}

	for _, v := range e.splitValues(value) {
		// Skip empty items so trailing/duplicate splitters (e.g. "a,,b" or
		// "a,") don't fail validation.
		if v == "" {
			continue
		}
		if err := e.validateType(v); err != nil {
			return err
		}
		if err := e.validateOptions(v, getRemote); err != nil {
			return err
		}
		if err := e.validateRegex(v); err != nil {
			return err
		}
	}
	return nil
}

// splitValues returns the individual values to validate. For a regular env var
// this is just [value]; for a multi-select env var the value is split on
// Splitter so each selected item can be validated against the same
// Type/Options/RemoteOptions/Regex constraints.
func (e *EnvVarSpec) splitValues(value string) []string {
	if !e.MultiSelect {
		return []string{value}
	}
	return strings.Split(value, e.GetSplitter())
}

func (e *EnvVarSpec) validateType(value string) error {
	if value == "" {
		return nil
	}
	switch e.Type {
	case "", "string", "password":
		return nil
	case "int":
		_, err := strconv.Atoi(value)
		return err
	case "quantity":
		if _, err := resource.ParseQuantity(value); err != nil {
			return fmt.Errorf("invalid quantity '%s': %w", value, err)
		}
		return nil
	case "bool":
		_, err := strconv.ParseBool(value)
		return err
	case "url":
		_, err := url.ParseRequestURI(value)
		return err
	case "ip":
		ip := net.ParseIP(value)
		if ip == nil {
			return fmt.Errorf("invalid ip '%s'", value)
		}
		return nil
	case "domain":
		errs := validation.IsDNS1123Subdomain(value)
		if len(errs) > 0 {
			return fmt.Errorf("invalid domain '%s'", value)
		}
		return nil
	case "email":
		_, err := mail.ParseAddress(value)
		if err != nil {
			return fmt.Errorf("invalid email '%s'", value)
		}
	}
	return nil
}

// validateOptions validates the given value against Options and/or RemoteOptions.
// getRemote lazily resolves (and caches) the remote option list so it is only
// fetched when a value cannot be satisfied by the local Options.
// Rules:
// - If both Options and RemoteOptions are set, value is valid if it is in either set.
// - If only Options is set, value must be in Options.
// - If only RemoteOptions is set, value must be in the fetched remote list.
// - If neither is set, any value is accepted.
func (e *EnvVarSpec) validateOptions(value string, getRemote func() ([]EnvValueOptionItem, error)) error {
	if value == "" {
		return nil
	}
	hasOptions := len(e.Options) > 0
	hasRemote := strings.TrimSpace(e.RemoteOptions) != ""

	if !hasOptions && !hasRemote {
		return nil
	}

	// Local options short-circuit, so we avoid a remote fetch when possible.
	if hasOptions && optionsContainValue(e.Options, value) {
		return nil
	}

	if !hasRemote {
		return fmt.Errorf("value not in options")
	}

	allowed, err := getRemote()
	if err != nil {
		return fmt.Errorf("invalid remoteOptions: %w", err)
	}
	if optionsContainValue(allowed, value) {
		return nil
	}
	if hasOptions {
		return fmt.Errorf("value not allowed by options or remoteOptions")
	}
	return fmt.Errorf("value not in remoteOptions")
}

func optionsContainValue(options []EnvValueOptionItem, v string) bool {
	for _, item := range options {
		if item.Value == v {
			return true
		}
	}
	return false
}

// fetchRemoteOptions fetches allowed values from a remote URL.
// Response body must be a JSON array of EnvValueOptionItem: [{"title":"A","value":"a"}, ...]
func fetchRemoteOptions(endpoint string) ([]EnvValueOptionItem, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse url failed: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}
	var items []EnvValueOptionItem
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode json failed: %w", err)
	}
	return items, nil
}

func (e *EnvVarSpec) validateRegex(value string) error {
	if e.Regex == "" {
		return nil
	}
	re, err := regexp2.Compile(e.Regex, regexp2.None)
	if err != nil {
		return fmt.Errorf("invalid regex: %w", err)
	}
	matched, matchErr := re.MatchString(value)
	if matchErr != nil {
		return fmt.Errorf("regex match error: %w", matchErr)
	}
	if !matched {
		return fmt.Errorf("value '%s' does not match regex '%s'", value, e.Regex)
	}
	return nil
}
