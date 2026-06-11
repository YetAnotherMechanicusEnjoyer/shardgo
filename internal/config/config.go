package config

import (
	"os"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Server struct {
		Addr     string `yaml:"addr"`
		IsMaster bool   `yaml:"is_master"`
	}
	Cluster struct {
		Nodes    []string `yaml:"nodes"`
		Replicas int      `yaml:"replicas"`
	}
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
