package config

import (
	"os"

	"github.com/charmbracelet/wishlist"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Listen string          `yaml:"listen"`
	Port   int64           `yaml:"port"`
	Users  []wishlist.User `yaml:"users"`
	Netbox *Netbox         `yaml:"netbox"`
}

type Netbox struct {
	Host         string `yaml:"host"`
	Token        string `yaml:"token"`
	FilterRole   string `yaml:"filter_role"`
	User         string `yaml:"user"`
	ForwardAgent bool   `yaml:"forward_agent"`
	OnlyActive   bool   `yaml:"only_active"`
	CAFile       string `yaml:"ca_file"`
	CAPemData    []byte `yaml:"-"`
	IgnoreTLS    bool   `yaml:"ignore_tls"`
}

const defaultConfigFile = "./config.yml"

func Get() (*Config, error) {
	f, err := os.ReadFile(defaultConfigFile)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return nil, err
	}
	if err := cfg.loadCAFile(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) loadCAFile() (err error) {
	if c.Netbox.CAFile == "" {
		return
	}
	c.Netbox.CAPemData, err = os.ReadFile(c.Netbox.CAFile)
	return
}
