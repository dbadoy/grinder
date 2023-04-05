package cft

import (
	"net"

	"github.com/dbadoy/grinder/pkg/checkpoint"
	"github.com/dbadoy/grinder/pkg/database"
)

var _ Engine = (*Solo)(nil)

type Solo struct {
	localnode net.Addr

	database.Database
	checkpoint.CheckpointHandler
}
