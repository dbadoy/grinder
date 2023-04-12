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

func (s *Server) loop() {}

// func (s *Server) readLoop() {}

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
	p, err := s.doHandkshake(conn)
	if err != nil {
		return
	}

	s.cohorts = append(s.cohorts, p)
}

func (s *Server) ClusterHash() []byte {
	h := sha256.New()
	h.Write([]byte(strings.Join(s.cfg.Cluster, "-")))
	return h.Sum(nil)[:20]
}
