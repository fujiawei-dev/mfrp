/*
 * @Date: 2022.01.22 19:21
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.22 19:21
 */

package main

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/fujiawei-dev/mfrp/pkg/models"
	"github.com/fujiawei-dev/mfrp/pkg/utils/conn"
	"github.com/fujiawei-dev/mfrp/pkg/utils/log"
)

const (
	sleepMinDuration = 1
	sleepMaxDuration = 60
)

var ErrAlreadyInUse = errors.New("already in use")

func ControlProcess(cli *models.ProxyClient, wait *sync.WaitGroup) {
	defer wait.Done()

	if cli.PassiveMode {
		runPassiveMode(cli)
	} else {
		for {
			if err := runActiveMode(cli); err != nil && err != ErrAlreadyInUse {
				panic(err)
			}

			time.Sleep(sleepMaxDuration)
		}
	}
}

func runActiveMode(cli *models.ProxyClient) (err error) {
	c, err := loginToServer(cli)

	if err != nil {
		log.Errorf("ProxyName [%s], connect to server failed!", cli.Name)
		return
	}

	defer c.Close()

	for {
		// ignore response content now
		if _, err = c.ReadLine(); err != nil {
			log.Errorf("ProxyName [%s], read from server error, %v", cli.Name, err)
			var sleep time.Duration = sleepMinDuration
			for {
				log.Debugf("ProxyName [%s], try to reconnect to server[%s:%d]...",
					cli.Name, conf.Common.ServerHost, conf.Common.ServerPort)

				if cRetry, errRetry := loginToServer(cli); errRetry == nil {
					c.Close()
					c = cRetry
					break
				}

				time.Sleep(time.Second * sleep)
				if sleep < sleepMaxDuration {
					sleep++
				}
			}

			continue
		}

		_ = cli.StartTunnel(conf.Common.ServerHost, conf.Common.ServerPort)
	}
}

func loginToServer(cli *models.ProxyClient) (c *conn.Conn, err error) {
	c = &conn.Conn{}

	if err = c.ConnectServer(conf.Common.ServerHost, conf.Common.ServerPort); err != nil {
		log.Errorf("ProxyName [%s], connect to server [%s:%d] error, %v",
			cli.Name, conf.Common.ServerHost, conf.Common.ServerPort, err)
		return
	}

	defer func() {
		if err != nil {
			c.Close()
		}
	}()

	req := &models.ClientCtlReq{
		Type:      models.ControlConn,
		ProxyName: cli.Name,
		Password:  cli.Password,
	}

	buf, _ := json.Marshal(req)
	if err = c.Write(string(buf) + "\n"); err != nil {
		log.Errorf("ProxyName [%s], write to server error, %v", cli.Name, err)
		return
	}

	res, err := c.ReadLine()
	if err != nil {
		log.Errorf("ProxyName [%s], read from server error, %v", cli.Name, err)
		return
	}

	log.Debugf("ProxyName [%s], read [%s]", cli.Name, res)

	clientCtlRes := &models.ClientCtlRes{}
	if err = json.Unmarshal([]byte(res), &clientCtlRes); err != nil {
		log.Errorf("ProxyName [%s], format server response error, %v", cli.Name, err)
		return
	}

	if clientCtlRes.Code != models.Success {
		if clientCtlRes.Code == models.AlreadyInUse {
			err = ErrAlreadyInUse
		} else {
			err = errors.New(clientCtlRes.Message)
		}
		log.Errorf("ProxyName [%s], start proxy error, %s", cli.Name, clientCtlRes.Message)
	}

	return
}

func runPassiveMode(cli *models.ProxyClient) {
	l, err := conn.Listen(cli.BindAddr, cli.BindPort)

	if err != nil {
		log.Errorf("create listener error, %v", err)
		os.Exit(-1)
	}

	log.Infof("start mfrpc %s on passive mode", l.Addr.String())

	for {
		remoteConn := l.GetConn()

		if localConn, err := cli.GetLocalConn(); err == nil {
			go conn.Join(localConn, remoteConn)
		} else {
			remoteConn.Close()
		}
	}
}
