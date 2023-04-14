package p2p

import (
	"crypto/sha256"
	"net"
	"strings"

	"github.com/dbadoy/grinder/pkg/checkpoint"
)

type Config struct {
	Cluster []string
}

type Server struct {
	listener *net.TCPListener

	localnode net.Addr
	cohorts   []*peer

	cp checkpoint.CheckpointHandler

	cfg *Config
}

func New(localnode *net.TCPAddr, cfg *Config) (*Server, error) {
	listener, err := net.ListenTCP("tcp", localnode)
	if err != nil {
		return nil, err
	}
	return &Server{listener: listener}, nil
}

func (s *Server) Run() {
	go s.loop()
	go s.acceptLoop()
}

func (s *Server) loop() {}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}

		if !strings.Contains("cluster list", conn.RemoteAddr().String()) {
			return
		}

		s.setupConn(conn)
	}
}

func (s *Server) setupConn(conn net.Conn) {
	peer, err := s.doHandkshake(conn)
	if err != nil {
		return
	}

	go peer.run()

	s.cohorts = append(s.cohorts, peer)
}

func (s *Server) broadcastMsg(msg msg) {
	for _, cohort := range s.cohorts {
		cohort.conn.writeMsg(msg)
	}
}

func ClusterHash(cluster []string) []byte {
	h := sha256.New()
	h.Write([]byte(strings.Join(cluster, "-")))
	return h.Sum(nil)[:20]
}
