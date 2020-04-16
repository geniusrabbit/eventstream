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
	mx sync.RWMutex

	LogLevel string `default:"debug" env:"LOG_LEVEL"`

	Stores  map[string]configItem `yaml:"stores" json:"stores"`
	Sources map[string]configItem `yaml:"sources" json:"sources"`
	Streams map[string]configItem `yaml:"streams" json:"streams"`

	Jaeger struct {
		AgentHost string `env:"JAEGER_AGENT_HOST"`
	}
}

func (cfg *config) String() string {
	data, _ := json.MarshalIndent(cfg, "", "  ")
	return string(data)
}

// Load eventstore config
func (cfg *config) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	file.Close()

	cfg.mx.Lock()
	defer cfg.mx.Unlock()

	switch strings.ToLower(filepath.Ext(filename)) {
	case ".yml", ".yaml":
		err = yaml.Unmarshal(data, cfg)
	case ".hcl":
		err = hcl.Unmarshal(data, cfg)
	case ".json":
		err = json.Unmarshal(data, cfg)
	default:
		err = errInvalidConfigFile
	}
	return err
}

// Validate config
func (cfg *config) Validate() error {
	if cfg == nil {
		return errInvalidConfig
	}
	if cfg.Stores == nil || len(cfg.Stores) < 1 {
		return errInvalidStoreConnect
	}
	if cfg.Sources == nil || len(cfg.Sources) < 1 {
		return errInvalidSourceConnect
	}
	if cfg.Streams == nil || len(cfg.Streams) < 1 {
		return errInvalidSourceStream
	}
	return nil
}

// IsDebug mode ON
func (cfg *config) IsDebug() bool {
	return strings.ToLower(cfg.LogLevel) == "debug"
}

// Config instance
var Config config
