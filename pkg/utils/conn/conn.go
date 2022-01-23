/*
 * @Date: 2022.01.22 17:56
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.22 17:56
 */

package conn

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/fujiawei-dev/mfrp/pkg/utils/log"
)

type Listener struct {
	Addr  net.Addr
	Conns chan *Conn
}

func (l *Listener) GetConn() (conn *Conn) {
	conn = <-l.Conns
	log.Infof("accept new tcp connection from %v", conn.GetRemoteAddr())
	return conn
}

type Conn struct {
	TcpConn *net.TCPConn
	Reader  *bufio.Reader
}

func (c *Conn) ConnectServer(host string, port int64) (err error) {
	serverAdder, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, serverAdder)
	if err != nil {
		return err
	}
	c.TcpConn = conn
	c.Reader = bufio.NewReader(c.TcpConn)
	return nil
}

func (c *Conn) GetRemoteAddr() (addr string) {
	return c.TcpConn.RemoteAddr().String()
}

func (c *Conn) GetLocalAddr() (addr string) {
	return c.TcpConn.LocalAddr().String()
}

func (c *Conn) ReadLine() (buff string, err error) {
	buff, err = c.Reader.ReadString('\n')
	return buff, err
}

func (c *Conn) Write(content string) (err error) {
	_, err = c.TcpConn.Write([]byte(content))
	return err
}

func (c *Conn) Close() {
	_ = c.TcpConn.Close()
}

func Listen(bindAddr string, bindPort int64) (l *Listener, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", bindAddr, bindPort))
	if err != nil {
		return l, err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return l, err
	}

	l = &Listener{
		Addr:  listener.Addr(),
		Conns: make(chan *Conn),
	}

	go func() {
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Errorf("accept new tcp connection error, %v", err)
				continue
			}

			c := &Conn{
				TcpConn: conn,
				Reader:  bufio.NewReader(conn),
			}

			l.Conns <- c
		}
	}()

	return l, err
}

func Join(c1 *Conn, c2 *Conn) {
	var wait sync.WaitGroup

	pipe := func(to *Conn, from *Conn) {
		defer wait.Done()

		_, _ = io.Copy(to.TcpConn, from.TcpConn)
	}

	wait.Add(2)

	go pipe(c1, c2)
	go pipe(c2, c1)

	log.Debugf(
		"join two conns, %s <-> %s <+> %s <-> %s",
		c2.GetRemoteAddr(),
		c2.GetLocalAddr(),
		c1.GetLocalAddr(),
		c1.GetRemoteAddr(),
	)

	wait.Wait()
}
