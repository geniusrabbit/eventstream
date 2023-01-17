package stream

import (
	"encoding/json"
	"reflect"

	"github.com/geniusrabbit/eventstream/internal/metrics"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	errInvalidStoreParameter  = errors.New("[storage::config] invalid the 'store' parameter")
	errInvalidSourceParameter = errors.New("[storage::config] invalid the 'source' parameter")
)

// Config of the stream
type Config struct {
	Name    string
	Debug   bool
	Store   string
	Source  string
	Where   string
	Metrics metrics.MetricList
	Raw     json.RawMessage
}

// Decode raw data to the target object
func (c *Config) Decode(v any) error {
	if len(c.Raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(c.Raw, v); err != nil {
		zap.L().Debug(`decode stream config`,
			zap.String(`json_raw`, string(c.Raw)),
			zap.String(`target_type`, reflect.TypeOf(v).String()),
			zap.Error(err))
		return errors.Wrap(err, `invalid decode stream config`)
	}
	return nil
}

// UnmarshalJSON data
func (c *Config) UnmarshalJSON(data []byte) (err error) {
	c.Raw = json.RawMessage(data)

	var conf struct {
		Store   string             `json:"store"`
		Source  string             `json:"source"`
		Where   string             `json:"where"`
		Metrics metrics.MetricList `json:"metrics"`
	}

	if err = json.Unmarshal(data, &conf); err == nil {
		c.Store = conf.Store
		c.Source = conf.Source
		c.Where = conf.Where
		c.Metrics = conf.Metrics
	}
	return err
}

// Validate config object
func (c *Config) Validate() error {
	if c.Store == "" {
		return errInvalidStoreParameter
	}
	if c.Source == "" {
		return errInvalidSourceParameter
	}
	return nil
}
