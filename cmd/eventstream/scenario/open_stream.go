package scenario

import (
	"fmt"

	cmdstorage "github.com/geniusrabbit/eventstream/cmd/eventstream/storage"
)

func streamOptsForStore(streamName, storeDriver string, item Item, debug bool) ([]any, error) {
	p := cmdstorage.StreamParams{
		Debug:        debug,
		SQLQuery:     itemString(item, "sql_query"),
		Target:       itemString(item, "target"),
		Fields:       item["fields"],
		IterateBy:    itemString(item, "iterate_by"),
		BufferSize:   itemUint(item, "buffer_size"),
		WriteTimeout: itemUint(item, "write_timeout"),
		Key:          itemString(item, "key"),
		Expiration:   itemDuration(item, "expiration"),
		Incrementor:  itemBool(item, "incrementor"),
		Method:       itemString(item, "method"),
		URL:          itemString(item, "url"),
		ContentType:  itemString(item, "content_type"),
	}
	if rawTargets, ok := item["targets"]; ok && rawTargets != nil {
		targets, err := parseStreamTargets(rawTargets)
		if err != nil {
			return nil, fmt.Errorf("stream %q targets: %w", streamName, err)
		}
		p.Targets = targets
	}
	return cmdstorage.StreamOpts(storeDriver, streamName, p)
}

func parseStreamTargets(raw any) ([]cmdstorage.Target, error) {
	list, ok := raw.([]any)
	if !ok {
		if m, ok := raw.(map[string]any); ok {
			list = []any{m}
		} else {
			return nil, fmt.Errorf("expected list, got %T", raw)
		}
	}
	out := make([]cmdstorage.Target, 0, len(list))
	for i, item := range list {
		m, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("[%d]: expected object, got %T", i, item)
		}
		fields, err := parseTargetFields(m["fields"])
		if err != nil {
			return nil, fmt.Errorf("[%d].fields: %w", i, err)
		}
		where, _ := m["where"].(string)
		out = append(out, cmdstorage.Target{Fields: fields, Where: where})
	}
	return out, nil
}

// parseTargetFields accepts a map (YAML/TOML/JSON) or an HCL1 nested block,
// which decodes as a singleton list of maps: fields { k = v } → [{k:v}].
func parseTargetFields(raw any) (map[string]any, error) {
	switch v := raw.(type) {
	case nil:
		return nil, nil
	case map[string]any:
		return v, nil
	case []any:
		out := make(map[string]any)
		for _, item := range v {
			m, ok := item.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("expected map entries, got %T", item)
			}
			for k, val := range m {
				out[k] = val
			}
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported type %T", raw)
	}
}
