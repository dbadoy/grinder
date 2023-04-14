package p2p

const (
	handshakeMsgType = byte(0x01) + iota
	ackMsgType
	pingMsgType
	pongMsgType
	checkpointMsgType
)

type msg interface {
	kind() byte
	clusterCode() []byte
	checkpoint() uint64
}

type (
	handshakeMsg struct {
		Version     int
		ClusterCode []byte
		Checkpoint  uint64
	}

	ackMsg struct {
		Version     int
		ClusterCode []byte
		Checkpoint  uint64
	}

	pingMsg struct {
		ClusterCode []byte
		Checkpoint  uint64
	}

	pongMsg struct {
		ClusterCode []byte
		Checkpoint  uint64
	}

	checkpointMsg struct {
		ClusterCode []byte
		Checkpoint  uint64
	}
)

func encodeMsg(m msg) ([]byte, error)

func decodeMsg(b []byte) (msg, error)

func (handshakeMsg) kind() byte             { return handshakeMsgType }
func (h *handshakeMsg) clusterCode() []byte { return h.ClusterCode }
func (h *handshakeMsg) checkpoint() uint64  { return h.Checkpoint }

func (ackMsg) kind() byte             { return ackMsgType }
func (a *ackMsg) clusterCode() []byte { return a.ClusterCode }
func (a *ackMsg) checkpoint() uint64  { return a.Checkpoint }

func (pingMsg) kind() byte             { return pingMsgType }
func (p *pingMsg) clusterCode() []byte { return p.ClusterCode }
func (p *pingMsg) checkpoint() uint64  { return p.Checkpoint }

func (pongMsg) kind() byte             { return pongMsgType }
func (p *pongMsg) clusterCode() []byte { return p.ClusterCode }
func (p *pongMsg) checkpoint() uint64  { return p.Checkpoint }

func (checkpointMsg) kind() byte             { return checkpointMsgType }
func (c *checkpointMsg) clusterCode() []byte { return c.ClusterCode }
func (c *checkpointMsg) checkpoint() uint64  { return c.Checkpoint }
