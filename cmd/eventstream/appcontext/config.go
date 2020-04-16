//
// @project geniusrabbit::eventstream 2017, 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2020
//

package appcontext

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hashicorp/hcl"

	yaml "gopkg.in/yaml.v2"
)

var (
	errInvalidConfig        = errors.New("[config] invalid config")
	errInvalidConfigFile    = errors.New("[config] invalid config file format")
	errInvalidStoreConnect  = errors.New("[config] invalid store config connect")
	errInvalidSourceConnect = errors.New("[config] invalid source config connect")
	errInvalidSourceStream  = errors.New("[config] invalid stream config connect")
)

type configItem map[string]interface{}

func (it configItem) Decode(v interface{}) error {
	raw, err := json.Marshal(it)
	if err != nil {
		return fmt.Errorf("[config] invalid item encoding: %s", err)
	}
	return json.Unmarshal(raw, v)
}

type config struct {
	mx      sync.RWMutex
	Debug   bool                  `yaml:"debug" json:"debug"`
	Stores  map[string]configItem `yaml:"stores" json:"stores"`
	Sources map[string]configItem `yaml:"sources" json:"sources"`
	Streams map[string]configItem `yaml:"streams" json:"streams"`
}

func (c *config) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}

// Load eventstore config
func (c *config) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	file.Close()

	c.mx.Lock()
	defer c.mx.Unlock()

	switch strings.ToLower(filepath.Ext(filename)) {
	case ".yml", ".yaml":
		err = yaml.Unmarshal(data, c)
	case ".hcl":
		err = hcl.Unmarshal(data, c)
	case ".json":
		err = json.Unmarshal(data, c)
	default:
		err = errInvalidConfigFile
	}
	return err
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
