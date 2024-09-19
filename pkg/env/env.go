package env

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Common   Common             `toml:"common"`
	Samples  Samples            `toml:"samples"`
	Services map[string]Service `toml:"services"`
}

type Common struct {
	Name string `toml:"name"`
}

type Samples struct {
	Prefix string   `toml:"prefix"`
	Images []string `toml:"images"`
}

type Service struct {
	Name    string `toml:"name"`
	Group   string `toml:"group"`
	Url     string `toml:"url"`
	Default bool   `toml:"default"`

	// Request can be two types - formdata or json
	RequestJSONTemplate  string `toml:"request_json_template"`
	RequestFormdataField string `toml:"request_formdata_field"`

	// Response can be only json
	ResponseXYField string `toml:"response_xy_field"`
}

func (s Service) String() string {
	return s.Name
}

type ServiceGroup string

// Parse TOML config
func ParseConfig(tomlConfig string) (*Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(tomlConfig, &conf); err != nil {
		return nil, fmt.Errorf("can't parse config: %w", err)
	}

	return &conf, nil
}

// Extract service groups from config
func (c *Config) GetServiceGroups() map[ServiceGroup][]Service {
	var groupedServices = make(map[ServiceGroup][]Service)
	for _, v := range c.Services {
		name := ServiceGroup(v.Group)
		groupedServices[name] = append(groupedServices[name], v)
	}

	return groupedServices
}

// Extract samples from config
func (c *Config) GetSamples() []string {
	var result []string
	for _, v := range c.Samples.Images {
		result = append(result, c.Samples.Prefix+v)
	}

	return result
}
