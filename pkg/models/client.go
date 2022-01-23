/*
 * @Date: 2022.01.22 18:14
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.22 18:14
 */

package models

import (
	"encoding/json"
	"github.com/fujiawei-dev/mfrp/pkg/utils/conn"
	"github.com/fujiawei-dev/mfrp/pkg/utils/log"
)

type ProxyClient struct {
	// Actively connect to the public server, double forwarding
	// self (nat) -> server1 (public) && self (nat) -> server2 (nat)  => client <-> server1 <-> server2
	Name      string `yaml:"Name,omitempty"`
	Password  string `yaml:"Password,omitempty"`
	LocalPort int64  `yaml:"LocalPort,omitempty"`

	// Passive listening port, then actively forwarding to local server
	// client -> self (public) -> server (nat) => client <-> server
	PassiveMode bool   `yaml:"PassiveMode,omitempty"`
	BindAddr    string `yaml:"BindAddr,omitempty"`
	BindPort    int64  `yaml:"BindPort,omitempty"`
}

func (p *ProxyClient) GetLocalConn() (c *conn.Conn, err error) {
	c = &conn.Conn{}
	if err = c.ConnectServer("127.0.0.1", p.LocalPort); err != nil {
		log.Errorf("ProxyName [%s], connect to local port error, %v", p.Name, err)
	}
	return
}

func (p *ProxyClient) GetRemoteConn(addr string, port int64) (c *conn.Conn, err error) {
	c = &conn.Conn{}

	defer func() {
		if err != nil {
			c.Close()
		}
	}()

	if err = c.ConnectServer(addr, port); err != nil {
		log.Errorf("ProxyName [%s], connect to server [%s:%d] error, %v", p.Name, addr, port, err)
		return
	}

	req := &ClientCtlReq{
		Type:      WorkConn,
		ProxyName: p.Name,
		Password:  p.Password,
	}

	buf, _ := json.Marshal(req)
	if err = c.Write(string(buf) + "\n"); err != nil {
		log.Errorf("ProxyName [%s], write to server error, %v", p.Name, err)
		return
	}

	return
}

func (p *ProxyClient) StartTunnel(serverAddr string, serverPort int64) (err error) {
	localConn, err := p.GetLocalConn()
	if err != nil {
		return
	}

	remoteConn, err := p.GetRemoteConn(serverAddr, serverPort)
	if err != nil {
		return
	}

	go conn.Join(localConn, remoteConn)

	return nil
}
