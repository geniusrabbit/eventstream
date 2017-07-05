//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package context

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl"

	yaml "gopkg.in/yaml.v2"
)

var (
	errInvalidConfig        = errors.New("Invalid config")
	errInvalidConfigFile    = errors.New("Invalid config file format")
	errInvalidStoreConnect  = errors.New("Invalid store config connect")
	errInvalidSourceConnect = errors.New("Invalid source config connect")
	errInvalidStreamTarget  = errors.New("Invalid stream target")
)

// StoreConfig description
type StoreConfig struct {
	Connect string  `yaml:"connect" json:"connect"`
	Options options `yaml:"options" json:"options"`
}

// Validate stream item
func (l StoreConfig) Validate() error {
	if "" == l.Connect {
		return errInvalidStoreConnect
	}
	return nil
}

// ConnectScheme name
func (l StoreConfig) ConnectScheme() string {
	return l.Connect[:strings.Index(l.Connect, "://")]
}

// SourceConfig description
type SourceConfig struct {
	Connect string  `yaml:"connect" json:"connect"`
	Format  string  `yaml:"format" json:"format"`
	Options options `yaml:"options" json:"options"`
}

// Validate stream item
func (l SourceConfig) Validate() error {
	if "" == l.Connect {
		return errInvalidSourceConnect
	}
	return nil
}

// ConnectScheme name
func (l SourceConfig) ConnectScheme() string {
	return l.Connect[:strings.Index(l.Connect, "://")]
}

// StreamConfig info
type StreamConfig struct {
	Store   string      `yaml:"store" json:"store"`
	Source  string      `yaml:"source" json:"source" default:"default"`
	RawItem string      `yaml:"rawitem" json:"rawitem"` // Depends from stream it could be SQL query or file raw
	Target  string      `yaml:"target" json:"target"`
	Fields  interface{} `yaml:"fields" json:"fields"`
	When    string      `yaml:"when" json:"when"`
	Options options     `yaml:"options" json:"options"`
}

// Validate log item
func (l StreamConfig) Validate() error {
	if "" == l.RawItem && "" == l.Target {
		return errInvalidStreamTarget
	}
	if "" == l.Source {
		l.Source = "default"
	}
	return nil
}

type config struct {
	Stores  map[string]StoreConfig  `yaml:"stores" json:"stores"`
	Sources map[string]SourceConfig `yaml:"sources" json:"sources"`
	Streams map[string]StreamConfig `yaml:"streams" json:"streams"`
}

// Load config
func (c *config) Load(filename string) error {
	file, err := os.Open(filename)
	if nil != err {
		return err
	}

	data, err := ioutil.ReadAll(file)
	if nil != err {
		return err
	}
	file.Close()

	switch strings.ToLower(filepath.Ext(filename)) {
	case ".yml", ".yaml":
		return yaml.Unmarshal(data, c)
	case ".hcl":
		return hcl.Unmarshal(data, c)
	}
	return errInvalidConfigFile
}

// Validate config
func (c *config) Validate() error {
	if nil == c || len(c.Stores) < 1 || len(c.Sources) < 1 || len(c.Streams) < 1 {
		return errInvalidConfig
	}

	for name, stream := range c.Streams {
		if err := stream.Validate(); nil != err {
			return fmt.Errorf("Stream [%s] %s", name, err.Error())
		}
	}

	for name, stream := range c.Streams {
		if err := stream.Validate(); nil != err {
			return fmt.Errorf("Stream [%s] %s", name, err.Error())
		} else if _, ok := c.Sources[stream.Source]; !ok {
			return fmt.Errorf("Stream [%s] Invalid source: %s", name, stream.Source)
		} else if _, ok := c.Stores[stream.Store]; !ok {
			return fmt.Errorf("Stream [%s] Invalid store: %s", name, stream.Store)
		}
	}
	return nil
}

// Config instance
var Config config
