package p2p

import (
	"net"
	"time"
)

const (
	maxPacketSize      = 128
	defaultReadTimeout = time.Second
)

type connection struct {
	conn net.Conn
}

func (c *connection) writeMsg(msg msg) error {
	b, err := encodeMsg(msg)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(b)
	return err
}

func (c *connection) readMsg() (msg, error) {
	b := make([]byte, maxPacketSize)
	n, err := c.conn.Read(b)
	if err != nil {
		return nil, err
	}
	return decodeMsg(b[:n])
}

func (c *connection) readMsgWithTimeout(timeout time.Duration) (msg, error) {
	c.conn.SetReadDeadline(time.Now().Add(timeout))
	return c.readMsg()
}

type peer struct {
	conn    *connection
	handler *handler

	closed chan struct{}
}

func (p *peer) run() {
	go p.readLoop()
	go p.pingLoop()
}

func (p *peer) readLoop() {
	for {
		msg, err := p.conn.readMsg()
		if err != nil {
			return
		}

		p.handler.handleMsg(p, msg)
	}
}

func (p *peer) pingLoop() {
	ping := time.NewTimer(0)
	defer ping.Stop()

	for {
		select {
		case <-ping.C:
			ping.Reset(0)

		case <-p.closed:
			return
		}
	}
}
