package stream

import (
	"encoding/json"
)

// Config of the storage
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
	return json.Unmarshal(c.Raw, v)
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
