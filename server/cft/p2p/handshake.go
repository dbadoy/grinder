package p2p

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"
)

var handshakeTimeout = time.Second

func (s *Server) doHandkshake(c net.Conn) (*peer, error) {
	conn := &connection{c}

	conn.writeMsg(&handshakeMsg{
		ClusterCode: s.ClusterHash(),
		Checkpoint:  s.cp.Checkpoint(),
	})

	msg, err := conn.readMsgWithTimeout(handshakeTimeout)
	if err != nil {
		return nil, err
	}

	ack := msg.(*ackMsg)
	if !bytes.Equal(s.ClusterHash(), ack.ClusterCode) {
		return nil, fmt.Errorf("has invalid cluster list: %s", c.RemoteAddr())
	}

	if ack.Checkpoint < s.cp.Checkpoint() {
		panic("into follower role")
	}

	// start checkpoint sync.
	if ack.Checkpoint > s.cp.Checkpoint() {
		conn.writeMsg(&checkpointMsg{
			ClusterCode: s.ClusterHash(),
			Checkpoint:  s.cp.Checkpoint(),
		})

		msg, err := conn.readMsgWithTimeout(handshakeTimeout)
		if err != nil {
			return nil, err
		}

		ack := msg.(*ackMsg)
		if s.cp.Checkpoint() != ack.Checkpoint {
			return nil, errors.New("invalid protocol")
		}
	}

	return &peer{
		conn:   conn,
		closed: make(chan struct{}),
	}, nil
}
