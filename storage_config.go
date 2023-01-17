package eventstream

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/geniusrabbit/eventstream/internal/utils"
)

var (
	errStorageEmptyConnection = errors.New("[storage] empty connection")
	errStorageUndefinedDriver = errors.New("[storage] undefined driver")
)

// StorageConfig of the storage
type StorageConfig struct {
	Debug   bool
	Connect string
	Driver  string
	Buffer  uint
	Raw     json.RawMessage
}

// Decode raw data to the target object
func (c *StorageConfig) Decode(v any) error {
	if err := json.Unmarshal(c.Raw, v); err != nil {
		return fmt.Errorf("decode storage config: %s", err.Error())
	}
	return nil
}

// UnmarshalJSON data
func (c *StorageConfig) UnmarshalJSON(data []byte) (err error) {
	var confData struct {
		Connect string `json:"connect"`
		Driver  string `json:"driver"`
		Buffer  uint   `json:"buffer"`
	}

	if err = json.Unmarshal(data, &confData); err != nil {
		return err
	}

	if confData.Buffer <= 0 {
		confData.Buffer = 1000
	}
	c.Connect = utils.PrepareValue(confData.Connect)
	c.Driver = confData.Driver
	c.Buffer = confData.Buffer
	c.Raw = json.RawMessage(data)

	if c.Driver == `` {
		if urlData := strings.Split(c.Connect, `://`); len(urlData) > 1 {
			c.Driver = urlData[0]
		}
	}
	return err
}

// Validate config
func (c *StorageConfig) Validate() error {
	if c.Connect == "" {
		return errStorageEmptyConnection
	}
	if c.Driver == "" {
		return errStorageUndefinedDriver
	}
	return nil
}
