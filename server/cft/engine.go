package cft

import (
	"net"

	"github.com/dbadoy/grinder/pkg/database"
)

type Engine interface {
	// Database
	Insert(key []byte, data database.Data) error
	Put(key []byte, data database.Data) error
	Delete(key []byte) error
	Exist(index string, key []byte) (bool, error)

	// Checkpoint
	Checkpoint() uint64
	SetCheckpoint(uint64) error
	Increase() error
	Decrease() error
}

type server interface {
	HasLeaderPermissions() bool

	LocalNode() net.Addr

	ClusterLeader() net.Addr

	CommitCheckpoint(uint64) error
}
