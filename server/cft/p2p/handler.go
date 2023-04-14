package p2p

import (
	"bytes"

	"github.com/dbadoy/grinder/pkg/checkpoint"
)

type handler struct {
	cp  checkpoint.CheckpointHandler
	cfg *Config
}

func (h *handler) validateMsg(m msg) bool {
	if m.checkpoint() != h.cp.Checkpoint() {
		return false
	}

	if bytes.Equal(m.clusterCode(), m.clusterCode()) {
		return false
	}

	return true
}

func (h *handler) handleMsg(peer *peer, msg msg) {
	if !h.validateMsg(msg) {
		return
	}

	switch msg.kind() {
	case handshakeMsgType:
		h.handleHandshake(peer, msg.(*handshakeMsg))

	case ackMsgType:

	case pingMsgType:
		h.handlePing(peer, msg.(*pingMsg))

	case pongMsgType:

	case checkpointMsgType:
		h.handleCheckpoint(peer, msg.(*checkpointMsg))
	}
}

func (h *handler) handleHandshake(peer *peer, hs *handshakeMsg) {
	peer.conn.writeMsg(&ackMsg{
		ClusterCode: ClusterHash(h.cfg.Cluster),
		Checkpoint:  h.cp.Checkpoint(),
	})
}

func (h *handler) handlePing(peer *peer, p *pingMsg) {
	peer.conn.writeMsg(&pongMsg{
		ClusterCode: ClusterHash(h.cfg.Cluster),
		Checkpoint:  h.cp.Checkpoint(),
	})
}

func (h *handler) handleCheckpoint(peer *peer, cp *checkpointMsg) {
	h.cp.SetCheckpoint(cp.Checkpoint)

	peer.conn.writeMsg(&ackMsg{
		ClusterCode: ClusterHash(h.cfg.Cluster),
		Checkpoint:  h.cp.Checkpoint(),
	})
}
