//
// @project geniusrabbit::eventstream 2017, 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2020
//

package appcontext

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/demdxx/goconfig"
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

type profilerConfig struct {
	Mode   string `json:"mode" yaml:"mode" default:"" env:"SERVER_PROFILE_MODE"`
	Listen string `json:"listen" yaml:"listen" cli:"profiler" default:"" env:"SERVER_PROFILE_LISTEN"`
}

// ConfigType contains all application options
type ConfigType struct {
	mx sync.RWMutex

	Config string `json:"-" cli:"config"`

	LogLevel string `default:"debug" env:"LOG_LEVEL"`

	Profile profilerConfig `yaml:"profile" json:"profile"`

	Stores  map[string]configItem `yaml:"stores" json:"stores"`
	Sources map[string]configItem `yaml:"sources" json:"sources"`
	Streams map[string]configItem `yaml:"streams" json:"streams"`

	Jaeger struct {
		AgentHost string `env:"JAEGER_AGENT_HOST"`
	}
}

func (cfg *ConfigType) String() string {
	data, _ := json.MarshalIndent(cfg, "", "  ")
	return string(data)
}

// ConfigFilepath name
func (cfg *ConfigType) ConfigFilepath() string {
	return cfg.Config
}

// Load eventstore config
func (cfg *ConfigType) Load() error {
	cfg.mx.Lock()
	defer cfg.mx.Unlock()
	return goconfig.Load(cfg)
}

// Validate config
func (cfg *ConfigType) Validate() error {
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
func (cfg *ConfigType) IsDebug() bool {
	return strings.ToLower(cfg.LogLevel) == `debug`
}

// Config instance
var Config ConfigType
