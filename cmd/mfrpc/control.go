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
	"io"
	"os"
	"sync"

	"github.com/fujiawei-dev/mfrp/pkg/models"
	"github.com/fujiawei-dev/mfrp/pkg/utils/conn"
	"github.com/fujiawei-dev/mfrp/pkg/utils/log"
)

func ControlProcess(cli *models.ProxyClient, wait *sync.WaitGroup) {
	defer wait.Done()

	if cli.PassiveMode {
		runPassiveMode(cli)
	} else {
		for {
			if err := runActiveMode(cli); err != nil {
				break
			}
		}
	}
}

func runActiveMode(cli *models.ProxyClient) (err error) {
	c := &conn.Conn{}

	if err = c.ConnectServer(conf.Common.ServerHost, conf.Common.ServerPort); err != nil {
		log.Errorf(
			"ProxyName [%s], connect to server [%s:%d] error, %v",
			cli.Name, conf.Common.ServerHost, conf.Common.ServerPort, err,
		)
		return
	}

	defer c.Close()

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
		log.Errorf("ProxyName [%s], start proxy error, %s", cli.Name, clientCtlRes.Message)

		if clientCtlRes.Code == models.AlreadyInUse {
			return nil
		}

		return errors.New(clientCtlRes.Message)
	}

	for {
		// ignore response content now
		if _, err = c.ReadLine(); err == io.EOF {
			log.Debugf("ProxyName [%s], server close this control conn", cli.Name)
			return nil
		} else if err != nil {
			log.Warnf("ProxyName [%s], read from server error, %v", cli.Name, err)
			return nil
		}

		_ = cli.StartTunnel(conf.Common.ServerHost, conf.Common.ServerPort)
	}
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
