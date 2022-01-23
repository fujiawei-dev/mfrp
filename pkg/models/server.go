/*
 * @Date: 2022.01.22 18:20
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.22 18:20
 */

package models

import (
	"container/list"
	"github.com/fujiawei-dev/mfrp/pkg/utils/log"
	"sync"

	"github.com/fujiawei-dev/mfrp/pkg/utils/conn"
)

const (
	Idle = iota
	Working
)

type ProxyServer struct {
	Name       string `yaml:"Name,omitempty"`
	Password   string `yaml:"Password,omitempty"`
	BindAddr   string `yaml:"BindAddr,omitempty"`
	ListenPort int64  `yaml:"ListenPort,omitempty"`

	Status       int64           `yaml:"-"`
	Listener     *conn.Listener  `yaml:"-"` // accept new connection from remote users
	CtlMsgChan   chan int64      `yaml:"-"` // every time accept a new user conn, put "1" to the channel
	CliConnChan  chan *conn.Conn `yaml:"-"` // get client conns from control goroutine
	UserConnList *list.List      `yaml:"-"` // store user conns
	Mutex        sync.Mutex      `yaml:"-"`
}

func (p *ProxyServer) Init() {
	p.Status = Idle
	p.CtlMsgChan = make(chan int64)
	p.CliConnChan = make(chan *conn.Conn)
	p.UserConnList = list.New()
}

func (p *ProxyServer) Lock() {
	p.Mutex.Lock()
}

func (p *ProxyServer) Unlock() {
	p.Mutex.Unlock()
}

func (p *ProxyServer) Close() {
	p.Lock()
	p.Status = Idle
	p.CtlMsgChan = make(chan int64)
	p.CliConnChan = make(chan *conn.Conn)
	p.UserConnList = list.New()
	p.Unlock()
}

func (p *ProxyServer) WaitUserConn() (res int64) {
	res = <-p.CtlMsgChan
	return
}

func (p *ProxyServer) Start() (err error) {
	if p.Listener == nil {
		p.Listener, err = conn.Listen(p.BindAddr, p.ListenPort)
		if err != nil {
			return err
		}
	}

	p.Status = Working

	// start a goroutine for listener
	go func() {
		for {
			// block
			c := p.Listener.GetConn()
			log.Debugf("ProxyName [%s], get one new user conn [%s]", p.Name, c.GetRemoteAddr())

			// put to list
			p.Lock()
			if p.Status != Working {
				log.Debugf("ProxyName [%s] is not working, new user conn close", p.Name)
				c.Close()
				p.Unlock()
				return
			}
			p.UserConnList.PushBack(c)
			p.Unlock()

			// put msg to control conn
			p.CtlMsgChan <- 1
		}
	}()

	// start another goroutine for join two conns from client and user
	go func() {
		for {
			cliConn := <-p.CliConnChan

			p.Lock()

			var userConn *conn.Conn

			if p.UserConnList.Len() > 0 {
				element := p.UserConnList.Front()

				if element != nil {
					userConn = element.Value.(*conn.Conn)
					p.UserConnList.Remove(element)

					// msg will transfer to another without modifying
					go conn.Join(cliConn, userConn)
				}
			}

			p.Unlock()
		}
	}()

	return nil
}
