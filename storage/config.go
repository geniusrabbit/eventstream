package storage

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrStreamEmptyConnection = errors.New("[storage] empty connection")
	ErrStreamUndefinedDriver = errors.New("[storage] undefined driver")
)

// Config of the storage
type Config struct {
	Debug   bool
	Connect string
	Driver  string
	Buffer  uint
	Raw     json.RawMessage
}

// Decode raw data to the target object
func (c *Config) Decode(v interface{}) error {
	if err := json.Unmarshal(c.Raw, v); err != nil {
		return fmt.Errorf("decode storage config: %s", err.Error())
	}
	return nil
}

// UnmarshalJSON data
func (c *Config) UnmarshalJSON(data []byte) (err error) {
	var confData struct {
		Connect string `json:"connect"`
		Driver  string `json:"driver"`
		Buffer  uint   `json:"buffer"`
	}

	if err = json.Unmarshal(data, &confData); err == nil {
		if confData.Buffer <= 0 {
			confData.Buffer = 100
		}
		c.Connect = confData.Connect
		c.Driver = confData.Driver
		c.Buffer = confData.Buffer
		c.Raw = json.RawMessage(data)
	}
	return
}

// Validate config
func (c *Config) Validate() error {
	if c.Connect == "" {
		return ErrStreamEmptyConnection
	}
	if c.Driver == "" {
		return ErrStreamUndefinedDriver
	}
	return nil
}
