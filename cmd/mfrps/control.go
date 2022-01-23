/*
 * @Date: 2022.01.22 19:04
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.22 19:04
 */

package main

import (
	"encoding/json"
	"fmt"

	"mfrp/pkg/models"
	"mfrp/pkg/utils/conn"
	"mfrp/pkg/utils/log"
)

func ProcessControlConn(l *conn.Listener) {
	for {
		// waits for mfrpc's connection
		c := l.GetConn()
		go controlWorker(c)
	}
}

// control connection from every client and server
func controlWorker(c *conn.Conn) {

	// the first message is from client to server
	// if error, close connection
	res, err := c.ReadLine()
	if err != nil {
		log.Warnf("read error, %v", err)
		return
	}

	log.Debugf("get: %s", res)

	clientCtlReq := &models.ClientCtlReq{}
	clientCtlRes := &models.ClientCtlRes{}

	if err = json.Unmarshal([]byte(res), &clientCtlReq); err != nil {
		log.Warnf("parse err: %v : %s", err, res)
		return
	}

	// check
	success, message, needRes, errorCode := checkProxy(clientCtlReq, c)

	if !success {
		clientCtlRes.Code = errorCode
		clientCtlRes.Message = message
	}

	if needRes {
		buf, _ := json.Marshal(clientCtlRes)
		err = c.Write(string(buf) + "\n")
		if err != nil {
			log.Warnf("write error, %v", err)
		}
	} else {
		// work conn, just return
		return
	}

	// only control conn needs to close
	defer c.Close()

	if !success {
		return
	}

	// others are from server to client
	server, ok := ProxyServers[clientCtlReq.ProxyName]
	if !ok {
		log.Warnf("ProxyName [%s] is not exist", clientCtlReq.ProxyName)
		return
	}

	serverCtlReq := &models.ClientCtlReq{}
	serverCtlReq.Type = models.WorkConn

	for {
		server.WaitUserConn()

		// notify mfrpc to start a working connection
		buf, _ := json.Marshal(serverCtlReq)
		if err = c.Write(string(buf) + "\n"); err != nil {
			log.Warnf("ProxyName [%s], write to mfrpc error, proxy exit", server.Name)
			server.Close()
			return
		}

		log.Debugf("ProxyName [%s], write to mfrpc to add work conn success", server.Name)
	}
}

func checkProxy(req *models.ClientCtlReq, c *conn.Conn) (success bool, message string, needRes bool, errorCode int) {
	success = false
	needRes = true

	// check if proxy name exist
	server, ok := ProxyServers[req.ProxyName]

	if !ok {
		message = fmt.Sprintf("ProxyName [%s] is not exist", req.ProxyName)
		errorCode = models.NotExists
		log.Warnf(message)
		return
	}

	// check password
	if req.Password != server.Password {
		message = fmt.Sprintf("ProxyName [%s], password is not correct", req.ProxyName)
		errorCode = models.InvalidPassword
		log.Warnf(message)
		return
	}

	// control conn
	if req.Type == models.ControlConn {
		if server.Status != models.Idle {
			// if no user connects but mfrpc disconnected, always in use
			server.CtlMsgChan <- 1

			if server.Status != models.Idle {
				message = fmt.Sprintf("ProxyName [%s], already in use", req.ProxyName)
				errorCode = models.AlreadyInUse
				log.Warnf(message)
				return
			}
		}

		// start proxy and listen for user conn, no block
		if err := server.Start(); err != nil {
			message = fmt.Sprintf("ProxyName [%s], start proxy error: %v", req.ProxyName, err.Error())
			errorCode = models.Unexpected
			log.Warnf(message)
			return
		}

		log.Infof("ProxyName [%s], start proxy success", req.ProxyName)
	} else if req.Type == models.WorkConn {
		// work conn
		needRes = false

		if server.Status != models.Working {
			log.Warnf("ProxyName [%s], is not working when it gets one new work conn", req.ProxyName)
			return
		}

		log.Infof("ProxyName [%s], gets one new work conn", req.ProxyName)

		server.CliConnChan <- c
	} else {
		message = fmt.Sprintf("ProxyName [%s], type [%d] unsupport", req.ProxyName, req.Type)
		errorCode = models.Unsupported
		log.Warnf(message)
		return
	}

	success = true
	return
}
