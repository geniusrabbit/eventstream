//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package context

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/geniusrabbit/eventstream"
	"github.com/hashicorp/hcl"

	yaml "gopkg.in/yaml.v2"
)

var (
	errInvalidConfig        = errors.New("Invalid config")
	errInvalidConfigFile    = errors.New("Invalid config file format")
	errInvalidStoreConnect  = errors.New("Invalid store config connect")
	errInvalidSourceConnect = errors.New("Invalid source config connect")
	errInvalidSourceStream  = errors.New("Invalid stream config connect")
)

type config struct {
	Stores  map[string]eventstream.ConfigItem `yaml:"stores" json:"stores"`
	Sources map[string]eventstream.ConfigItem `yaml:"sources" json:"sources"`
	Streams map[string]eventstream.ConfigItem `yaml:"streams" json:"streams"`
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
	if c == nil {
		return errInvalidConfig
	}
	if c.Stores == nil || len(c.Stores) < 1 {
		return errInvalidStoreConnect
	}
	if c.Sources == nil || len(c.Sources) < 1 {
		return errInvalidSourceConnect
	}
	if c.Streams == nil || len(c.Streams) < 1 {
		return errInvalidSourceStream
	}
	return nil
}

// Config instance
var Config config
