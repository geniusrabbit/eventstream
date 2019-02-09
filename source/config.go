package source

import "encoding/json"

// Config of the source connection
type Config struct {
	Debug   bool
	Connect string
	Driver  string
	Format  string
	Raw     json.RawMessage
}

// Decode raw data to the target object
func (c *Config) Decode(v interface{}) error {
	return json.Unmarshal(c.Raw, v)
}

// UnmarshalJSON data
func (c *Config) UnmarshalJSON(data []byte) (err error) {
	var confData struct {
		Connect string `json:"connect"`
		Driver  string `json:"driver"`
		Format  string `json:"format"`
	}

	if err = json.Unmarshal(data, &confData); err == nil {
		c.Connect = confData.Connect
		c.Driver = confData.Driver
		c.Format = confData.Format
		c.Raw = json.RawMessage(data)
	}
	return
}
