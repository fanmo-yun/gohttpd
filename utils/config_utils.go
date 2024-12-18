package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig    `yaml:"server"`
	Static  HtmlConfig      `yaml:"html"`
	Logger  LoggerConfig    `yaml:"logger"`
	Gzip    bool            `yaml:"gzip"`
	Custom  []CustomConfig  `yaml:"custom"`
	Proxy   []ProxyConfig   `yaml:"proxy"`
	Backend []BackendConfig `yaml:"backend"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type HtmlConfig struct {
	Dirpath string `yaml:"path"`
	Index   string `yaml:"index"`
	Try     bool
}

type LoggerConfig struct {
	Out   string `yaml:"out"`
	Level string `yaml:"level"`
}

type CustomConfig struct {
	Urlpath  string `yaml:"url"`
	Filepath string `yaml:"file"`
}

type ProxyConfig struct {
	PathPrefix string `yaml:"prefix"`
	TargetURL  string `yaml:"target"`
}

type BackendConfig struct {
	BackendURL string `yaml:"url"`
}

func DefaultServer() *ServerConfig {
	return &ServerConfig{
		Host: "0.0.0.0",
		Port: "80",
	}
}

func DefaultHtml() *HtmlConfig {
	return &HtmlConfig{
		Dirpath: "html",
		Index:   "index.html",
		Try:     false,
	}
}

func DefaultLogger() *LoggerConfig {
	return &LoggerConfig{
		Out:   "stdout",
		Level: "info",
	}
}

func TrimSpace(in string) string { return strings.Trim(in, " ") }

func (c *Config) CoverConfig() error {
	if reflect.DeepEqual(c.Server, ServerConfig{}) {
		c.Server = *DefaultServer()
	} else {
		if TrimSpace(c.Server.Host) == "" {
			c.Server.Host = "0.0.0.0"
		} else if TrimSpace(c.Server.Port) == "" {
			c.Server.Port = "80"
		}
	}

	if reflect.DeepEqual(c.Logger, LoggerConfig{}) {
		c.Logger = *DefaultLogger()
	} else {
		if TrimSpace(c.Logger.Out) == "" {
			c.Logger.Out = "console"
		} else if TrimSpace(c.Logger.Level) == "" {
			c.Logger.Level = "info"
		}
	}

	if reflect.DeepEqual(c.Static, HtmlConfig{}) {
		c.Static = *DefaultHtml()
	} else {
		if TrimSpace(c.Static.Dirpath) == "" {
			c.Static.Dirpath = "html"
		} else if TrimSpace(c.Static.Index) == "" {
			c.Static.Index = "index.html"
		} else {
			s := strings.Split(c.Static.Index, " ")
			if strings.ToLower(s[0]) == "try" {
				c.Static.Try = true
				if len(s) > 1 {
					c.Static.Index = s[1]
				} else {
					return fmt.Errorf("html config Cannot Init")
				}
			}
		}
	}
	return nil
}

func LoadConfig() (*Config, error) {
	var config Config
	configPath := filepath.Join("conf", "gohttpd.yaml")
	confData, readErr := os.ReadFile(configPath)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read config file: %v", readErr)
	}
	if unmarshalErr := yaml.Unmarshal(confData, &config); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %v", unmarshalErr)
	}
	if err := config.CoverConfig(); err != nil {
		return nil, err
	}
	return &config, nil
}
