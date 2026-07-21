// Package scenario loads pipeline definitions (stores, sources, streams)
// from YAML, TOML, HCL, or JSON files.
package scenario

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrEmptyStores  = errors.New("[scenario] stores are required")
	ErrEmptySources = errors.New("[scenario] sources are required")
	ErrEmptyStreams = errors.New("[scenario] streams are required")
)

// Item is a raw named config block decoded into typed configs via JSON.
type Item map[string]any

// Decode marshals the raw item into the target typed config.
func (it Item) Decode(v any) error {
	raw, err := json.Marshal(it)
	if err != nil {
		return fmt.Errorf("[scenario] invalid item encoding: %w", err)
	}
	return json.Unmarshal(raw, v)
}

// Scenario describes a complete eventstream pipeline configuration.
type Scenario struct {
	Version     string          `json:"version" yaml:"version" toml:"version" hcl:"version"`
	Description string          `json:"description" yaml:"description" toml:"description" hcl:"description"`
	Stores      map[string]Item `json:"stores" yaml:"stores" toml:"stores" hcl:"stores"`
	Sources     map[string]Item `json:"sources" yaml:"sources" toml:"sources" hcl:"sources"`
	Streams     map[string]Item `json:"streams" yaml:"streams" toml:"streams" hcl:"streams"`
}

// Validate ensures the scenario has the required sections.
func (s *Scenario) Validate() error {
	if s == nil {
		return errors.New("[scenario] nil scenario")
	}
	if len(s.Stores) < 1 {
		return ErrEmptyStores
	}
	if len(s.Sources) < 1 {
		return ErrEmptySources
	}
	if len(s.Streams) < 1 {
		return ErrEmptyStreams
	}
	return nil
}
