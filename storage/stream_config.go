package storage

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	errInvalidStoreParameter  = errors.New("[storage::config] invalid the 'store' parameter")
	errInvalidSourceParameter = errors.New("[storage::config] invalid the 'source' parameter")
)

// StreamConfig of the stream
type StreamConfig struct {
	Name   string
	Debug  bool
	Store  string
	Source string
	Where  string
	Raw    json.RawMessage
}

// Decode raw data to the target object
func (c *StreamConfig) Decode(v interface{}) error {
	if err := json.Unmarshal(c.Raw, v); err != nil {
		return fmt.Errorf("decode stream config: %s", err.Error())
	}
	return nil
}

// UnmarshalJSON data
func (c *StreamConfig) UnmarshalJSON(data []byte) (err error) {
	c.Raw = json.RawMessage(data)

	var conf struct {
		Store  string `json:"store"`
		Source string `json:"source"`
		Where  string `json:"where"`
	}

	if err = json.Unmarshal(data, &conf); err == nil {
		c.Store = conf.Store
		c.Source = conf.Source
		c.Where = conf.Where
	}

	return err
}

// Validate config object
func (c *StreamConfig) Validate() error {
	if c.Store == "" {
		return errInvalidStoreParameter
	}
	if c.Source == "" {
		return errInvalidSourceParameter
	}
	return nil
}
