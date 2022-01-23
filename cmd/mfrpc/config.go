/*
 * @Date: 2022.01.22 19:15
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.22 19:15
 */

package main

import (
	"io/fs"
	"io/ioutil"

	"gopkg.in/yaml.v3"
	"mfrp/pkg/models"
)

var ProxyClients = make(map[string]*models.ProxyClient)

func InitProxyClients() {
	for i := range conf.ProxyClients {
		proxyClient := conf.ProxyClients[i]
		ProxyClients[proxyClient.Name] = proxyClient
	}
}

var conf = &Config{
	Common: Common{
		ServerHost: "localhost",
		ServerPort: 9527,
		LogLevel:   "debug",
	},
	ProxyClients: []*models.ProxyClient{
		{
			PassiveMode: true,
			BindAddr:    "0.0.0.0",
			BindPort:    18124,
			LocalPort:   8080,
		},
		{
			Name:      "mfrp",
			Password:  "mfrp",
			LocalPort: 8080,
		},
	},
}

type Config struct {
	Common       Common                `yaml:"Common,omitempty"`
	ProxyClients []*models.ProxyClient `yaml:"ProxyServers,omitempty"`
}

type Common struct {
	ServerHost string `yaml:"ServerHost,omitempty"`
	ServerPort int64  `yaml:"ServerPort,omitempty"`
	LogLevel   string `yaml:"LogLevel,omitempty"`
}

func (c *Config) Load(path string) error {
	if buf, err := ioutil.ReadFile(path); err == nil {
		return yaml.Unmarshal(buf, c)
	}

	if buf, err := yaml.Marshal(c); err == nil {
		return ioutil.WriteFile(path, buf, fs.ModePerm)
	}

	return nil
}
