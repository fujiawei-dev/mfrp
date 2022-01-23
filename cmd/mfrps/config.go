/*
 * @Date: 2022.01.22 18:30
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.22 18:30
 */

package main

import (
	"io/fs"
	"io/ioutil"

	"github.com/fujiawei-dev/mfrp/pkg/models"
	"gopkg.in/yaml.v3"
)

var ProxyServers = make(map[string]*models.ProxyServer)

func InitProxyServers() {
	for i := range conf.ProxyServers {
		proxyServer := conf.ProxyServers[i]
		proxyServer.Init()

		ProxyServers[proxyServer.Name] = proxyServer
	}
}

var conf = &Config{
	Common: Common{
		BindAddr: "0.0.0.0",
		BindPort: 9527,
		LogLevel: "debug",
	},
	ProxyServers: []*models.ProxyServer{
		{
			Name:       "mfrp",
			Password:   "mfrp",
			BindAddr:   "0.0.0.0",
			ListenPort: 18123,
		},
	},
}

type Config struct {
	Common       Common                `yaml:"Common,omitempty"`
	ProxyServers []*models.ProxyServer `yaml:"ProxyServers,omitempty"`
}

type Common struct {
	BindAddr string `yaml:"BindAddr,omitempty"`
	BindPort int64  `yaml:"BindPort,omitempty"`
	LogLevel string `yaml:"LogLevel,omitempty"`
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
