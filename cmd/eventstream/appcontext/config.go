// Package appcontext holds per-command runtime configuration for the
// eventstream service binary (logging, profiling, service name).
package appcontext

import (
	"encoding/json"
	"strings"
)

type profilerConfig struct {
	Mode   string `json:"mode" yaml:"mode" toml:"mode" default:"" env:"SERVER_PROFILE_MODE"`
	Listen string `json:"listen" yaml:"listen" toml:"listen" cli:"profiler" default:"" env:"SERVER_PROFILE_LISTEN"`
}

// ConfigType contains application runtime options (not the pipeline scenario).
// Each command loads its own instance via appcmd / goconfig — there is no process-global config.
type ConfigType struct {
	ConfigFile  string `json:"-" cli:"config" env:"CONFIG"`
	ServiceName string `json:"service_name" yaml:"service_name" toml:"service_name" env:"SERVICE_NAME" default:"eventstream"`

	LogAddr    string `json:"log_addr" yaml:"log_addr" toml:"log_addr" default:"" env:"LOG_ADDR"`
	LogLevel   string `json:"log_level" yaml:"log_level" toml:"log_level" default:"debug" env:"LOG_LEVEL" cli:"log-level"`
	LogEncoder string `json:"log_encoder" yaml:"log_encoder" toml:"log_encoder" default:"" env:"LOG_ENCODER"`

	Profile profilerConfig `json:"profile" yaml:"profile" toml:"profile"`
}

func (cfg *ConfigType) String() string {
	data, _ := json.MarshalIndent(cfg, "", "  ")
	return string(data)
}

// IsDebug mode ON
func (cfg *ConfigType) IsDebug() bool {
	return strings.EqualFold(cfg.LogLevel, `debug`)
}
