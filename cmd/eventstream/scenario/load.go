package scenario

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/hashicorp/hcl"
	"gopkg.in/yaml.v3"
)

// LoadFile loads a scenario from a YAML, TOML, HCL, or JSON file.
func LoadFile(path string) (*Scenario, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("[scenario] read %s: %w", path, err)
	}
	return Load(data, filepath.Ext(path))
}

// Load parses scenario bytes using the format implied by ext
// (e.g. ".yml", ".yaml", ".toml", ".hcl", ".json").
func Load(data []byte, ext string) (*Scenario, error) {
	var sc Scenario
	ext = strings.ToLower(ext)
	var err error
	switch ext {
	case ".yml", ".yaml":
		err = yaml.Unmarshal(data, &sc)
	case ".toml":
		err = toml.Unmarshal(data, &sc)
	case ".json":
		err = json.Unmarshal(data, &sc)
	case ".hcl":
		err = unmarshalHCL(data, &sc)
	default:
		return nil, fmt.Errorf("[scenario] unsupported config format: %s", ext)
	}
	if err != nil {
		return nil, fmt.Errorf("[scenario] %s: %w", strings.TrimPrefix(ext, "."), err)
	}
	normalizeMaps(&sc)
	return &sc, nil
}

func unmarshalHCL(data []byte, sc *Scenario) error {
	var root any
	if err := hcl.Unmarshal(data, &root); err != nil {
		return err
	}
	// Re-encode through JSON so nested HCL types become []any / map[string]any.
	raw, err := json.Marshal(root)
	if err != nil {
		return err
	}
	var tree any
	if err := json.Unmarshal(raw, &tree); err != nil {
		return err
	}
	tree = normalizeHCLTree(tree)
	raw, err = json.Marshal(tree)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, sc)
}

// normalizeHCLTree converts hashicorp/hcl nested block shapes into plain maps.
//
// Shape A (one block):
//
//	"sources": [ { "nats_1": [ { ... } ] } ]
//
// Shape B (many blocks merged):
//
//	"stores": [ { "ping": [ {...} ], "nats_store": [ {...} ] } ]
//
// Both become:
//
//	{ "nats_1": { ... } } / { "ping": {...}, "nats_store": {...} }
func normalizeHCLTree(v any) any {
	switch t := v.(type) {
	case []any:
		if m, ok := sectionListToMap(t); ok {
			return m
		}
		out := make([]any, len(t))
		for i, item := range t {
			out[i] = normalizeHCLTree(item)
		}
		return out
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, item := range t {
			out[k] = normalizeHCLTree(item)
		}
		return out
	default:
		return v
	}
}

func sectionListToMap(list []any) (map[string]any, bool) {
	if len(list) == 0 {
		return nil, false
	}

	// Shape B: single element that is a multi-key named-block map.
	if len(list) == 1 {
		if m, ok := list[0].(map[string]any); ok && looksLikeNamedBlocks(m) {
			return unwrapNamedBlocks(m), true
		}
	}

	// Shape A: list of single-key maps.
	out := make(map[string]any, len(list))
	for _, item := range list {
		m, ok := item.(map[string]any)
		if !ok || len(m) != 1 {
			return nil, false
		}
		for name, body := range m {
			if _, exists := out[name]; exists {
				return nil, false
			}
			out[name] = unwrapSingleton(body)
		}
	}
	return out, true
}

func looksLikeNamedBlocks(m map[string]any) bool {
	if len(m) == 0 {
		return false
	}
	for _, body := range m {
		arr, ok := body.([]any)
		if !ok || len(arr) != 1 {
			return false
		}
		if _, ok := arr[0].(map[string]any); !ok {
			return false
		}
	}
	return true
}

func unwrapNamedBlocks(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for name, body := range m {
		out[name] = unwrapSingleton(body)
	}
	return out
}

func unwrapSingleton(v any) any {
	switch t := v.(type) {
	case []any:
		if len(t) == 1 {
			return normalizeHCLTree(t[0])
		}
		return normalizeHCLTree(t)
	default:
		return normalizeHCLTree(v)
	}
}

func normalizeMaps(sc *Scenario) {
	if sc.Stores == nil {
		sc.Stores = map[string]Item{}
	}
	if sc.Sources == nil {
		sc.Sources = map[string]Item{}
	}
	if sc.Streams == nil {
		sc.Streams = map[string]Item{}
	}
}
