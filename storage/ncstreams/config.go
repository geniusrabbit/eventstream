package ncstreams

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

var (
	errInvalidMetricsItemConfig = errors.New("[metrics] invalid metrics item config")
)

type configFields map[string]any

func (cf *configFields) UnmarshalJSON(data []byte) (err error) {
	if bytes.HasPrefix(data, []byte("[")) {
		var manyf []map[string]any
		if err = json.Unmarshal(data, &manyf); err == nil {
			*cf = map[string]any{}
			for _, fields := range manyf {
				for key, val := range fields {
					(*cf)[key] = val
				}
			}
		}
	} else {
		var fields map[string]any
		if err = json.Unmarshal(data, &fields); err == nil {
			*cf = fields
		}
	}
	return nil
}

type configItem struct {
	Fields configFields `json:"fields"`
	Where  string       `json:"where"`
}

type config struct {
	Targets []configItem `json:"targets"`
}
