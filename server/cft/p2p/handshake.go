package p2p

import (
	"bytes"
	"errors"
	"fmt"
	"net"
)

func (s *Server) doHandkshake(c net.Conn) (*peer, error) {
	conn := &connection{c}

	conn.writeMsg(&handshakeMsg{
		ClusterCode: ClusterHash(s.cfg.Cluster),
		Checkpoint:  s.cp.Checkpoint(),
	})

	msg, err := conn.readMsgWithTimeout(defaultReadTimeout)
	if err != nil {
		return nil, err
	}

	ack := msg.(*ackMsg)

	if !bytes.Equal(ClusterHash(s.cfg.Cluster), ack.ClusterCode) {
		return nil, fmt.Errorf("has invalid cluster list: %s", c.RemoteAddr())
	}

	if ack.Checkpoint < s.cp.Checkpoint() {
		panic("into follower role")
	}

	// start checkpoint sync.
	if ack.Checkpoint > s.cp.Checkpoint() {
		conn.writeMsg(&checkpointMsg{
			ClusterCode: ClusterHash(s.cfg.Cluster),
			Checkpoint:  s.cp.Checkpoint(),
		})

		msg, err := conn.readMsgWithTimeout(defaultReadTimeout)
		if err != nil {
			return nil, err
		}

		ack := msg.(*ackMsg)
		if s.cp.Checkpoint() != ack.Checkpoint {
			return nil, errors.New("invalid protocol")
		}
	}

	return &peer{
		conn: conn,
		handler: &handler{
			cp:  s.cp,
			cfg: s.cfg,
		},
		closed: make(chan struct{}),
	}, nil
}
