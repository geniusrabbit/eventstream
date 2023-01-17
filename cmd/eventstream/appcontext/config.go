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
	errInvalidStoreConnect  = errors.New("[config] invalid store config connect")
	errInvalidSourceConnect = errors.New("[config] invalid source config connect")
	errInvalidSourceStream  = errors.New("[config] invalid stream config connect")
)

type configItem map[string]any

func (it configItem) Decode(v any) error {
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

	Config      string `json:"-" cli:"config"`
	ServiceName string `json:"service_name" yaml:"service_name" toml:"service_name" env:"SERVICE_NAME" default:"eventstream"`

	LogAddr    string `json:"log_addr" yaml:"log_addr" toml:"log_addr" default:"" env:"LOG_ADDR"`
	LogLevel   string `json:"log_level" yaml:"log_level" toml:"log_level" default:"debug" env:"LOG_LEVEL"`
	LogEncoder string `json:"log_encoder" yaml:"log_encoder" toml:"server" env:"LOG_ENCODER"`

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
	return strings.EqualFold(cfg.LogLevel, `debug`)
}

// Config instance
var Config ConfigType
