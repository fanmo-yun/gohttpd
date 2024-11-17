package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"go.uber.org/zap"
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
	Etcd    EtcdConfig      `yaml:"etcd"`
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

type EtcdConfig struct {
	Endpoints   []string `yaml:"endpoints"`
	ServiceName string   `yaml:"servicename"`
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

func (c *Config) CoverConfig() {
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
					zap.L().Fatal("gohttpd: html config Cannot Init")
				}
			}
		}
	}

	if reflect.DeepEqual(c.Etcd, EtcdConfig{}) || len(c.Etcd.Endpoints) == 0 || TrimSpace(c.Etcd.ServiceName) == "" {
		zap.L().Fatal("gohttpd: etcd config required")
	}
}

func LoadConfig() *Config {
	var config Config
	configPath := filepath.Join("conf", "gohttpd.yaml")
	confData, readErr := os.ReadFile(configPath)
	if readErr != nil {
		zap.L().Fatal("gohttpd: Config Cannot Init", zap.String("config", readErr.Error()))
	}
	if unmarshalErr := yaml.Unmarshal(confData, &config); unmarshalErr != nil {
		zap.L().Fatal("gohttpd: Config Cannot Unmarshal", zap.String("config", unmarshalErr.Error()))
	}
	config.CoverConfig()
	return &config
}
