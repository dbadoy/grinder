package p2p

import (
	"errors"
	"net"
	"time"
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
	b := make([]byte, 128)
	n, err := c.conn.Read(b)
	if err != nil {
		return nil, err
	}
	return decodeMsg(b[:n])
}

func (c *connection) readMsgWithTimeout(timeout time.Duration) (msg, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	recv := make(chan msg, 1)
	go func() {
		msg, err := c.readMsg()
		if err != nil {
			return
		}
		recv <- msg
	}()

	for {
		select {
		case msg := <-recv:
			return msg, nil
		case <-timer.C:
			return nil, errors.New("timeout")
		}
	}
}

type peer struct {
	conn   *connection
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

		_ = msg
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
