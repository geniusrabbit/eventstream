package stream

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	errInvalidStoreParameter  = errors.New("[storage::config] invalid the 'store' parameter")
	errInvalidSourceParameter = errors.New("[storage::config] invalid the 'source' parameter")
)

// Config of the stream
type Config struct {
	Name   string
	Debug  bool
	Store  string
	Source string
	Where  string
	Raw    json.RawMessage
}

// Decode raw data to the target object
func (c *Config) Decode(v interface{}) error {
	if len(c.Raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(c.Raw, v); err != nil {
		return fmt.Errorf("decode stream config: %s", err.Error())
	}
	return nil
}

// UnmarshalJSON data
func (c *Config) UnmarshalJSON(data []byte) (err error) {
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
func (c *Config) Validate() error {
	if c.Store == "" {
		return errInvalidStoreParameter
	}
	if c.Source == "" {
		return errInvalidSourceParameter
	}
	return nil
}
