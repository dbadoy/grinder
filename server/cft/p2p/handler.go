package p2p

func (s *Server) handleCheckpoint(peer *peer, cp *checkpointMsg) {
	if s.cp.Checkpoint() >= cp.Checkpoint {
		//
		return
	}

	s.cp.SetCheckpoint(cp.Checkpoint)

	peer.conn.writeMsg(&ackMsg{
		ClusterCode: s.ClusterHash(),
		Checkpoint:  s.cp.Checkpoint(),
	})
}

func (s *Server) handleHandshake(peer *peer, hs *handshakeMsg) {
	peer.conn.writeMsg(&ackMsg{
		ClusterCode: s.ClusterHash(),
		Checkpoint:  s.cp.Checkpoint(),
	})
}
