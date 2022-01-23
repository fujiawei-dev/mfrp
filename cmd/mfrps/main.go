/*
 * @Date: 2022.01.22 19:13
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.22 19:13
 */

package main

import (
	"os"

	"mfrp/pkg/utils/conn"
	"mfrp/pkg/utils/log"
)

func main() {
	if err := conf.Load("./mfrps.yaml"); err != nil {
		os.Exit(-1)
	}

	InitProxyServers()

	l, err := conn.Listen(conf.Common.BindAddr, conf.Common.BindPort)

	if err != nil {
		log.Errorf("create listener error, %v", err)
		os.Exit(-1)
	}

	log.Infof("mfrps %s waiting for connection", l.Addr.String())

	ProcessControlConn(l)
}
