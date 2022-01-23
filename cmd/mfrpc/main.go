/*
 * @Date: 2022.01.22 19:26
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.22 19:26
 */

package main

import (
	"os"
	"sync"

	"mfrp/pkg/utils/log"
)

func main() {
	if err := conf.Load("./mfrpc.yaml"); err != nil {
		os.Exit(-1)
	}

	InitProxyClients()

	// wait until all control goroutine exit
	var wait sync.WaitGroup

	wait.Add(len(ProxyClients))

	for _, client := range ProxyClients {
		go ControlProcess(client, &wait)
	}

	log.Infof("%d mfrpc waiting for connection", len(ProxyClients))

	wait.Wait()

	log.Warnf("all proxy exit!")
}
