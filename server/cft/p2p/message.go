package p2p

const (
	handshakeMsgType = byte(0x01) + iota
	pingMsgType
	pongMsgType
	syncMsgType
	checkpointMsgType
)

type msg interface {
	// kind() byte
}

type (
	handshakeMsg struct {
		ClusterCode []byte
		Checkpoint  uint64
	}

	ackMsg struct {
		ClusterCode []byte
		Checkpoint  uint64
	}

	pingMsg struct{}

	pongMsg struct{}

	syncMsg struct{}

	checkpointMsg struct {
		ClusterCode []byte
		Checkpoint  uint64
	}
)

func encodeMsg(m msg) ([]byte, error)

func decodeMsg(b []byte) (msg, error)
