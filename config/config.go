package config

import (
	"sync"
)

var (
	allconfig *AllConfig
	onceAll   sync.Once
)

// AllConfig is the configuration for a `mc all` command
type AllConfig struct {
	Path           string `yaml:"path"`
	InputFormat    string `yaml:"input-format"`
	OutputFormat   string `yaml:"output-format"`
	WithProcessBar bool   `yaml:"with-process-bar"`
	WorkPoolSize   uint   `yaml:"work-pool-size"`
	RemoveSourses  bool   `yaml:"remove-sourses"`
}

func init() {
	onceAll.Do(
		func() {
			allconfig = &AllConfig{}
		},
	)
}

// GetAllConfig will return the config sigleton
func GetAllConfig() *AllConfig {
	return allconfig
}
