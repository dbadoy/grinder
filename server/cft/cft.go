package cft

import (
	"net"

	"github.com/dbadoy/grinder/pkg/checkpoint"
	"github.com/dbadoy/grinder/pkg/database"
)

// var _ Engine = (*CFT)(nil)

type CFT struct {
	localnode net.Addr
	srv       server
	db        database.Database
	cp        checkpoint.CheckpointHandler
}
