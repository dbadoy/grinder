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

func NewSoloEngine(local net.Addr, db database.Database, cp checkpoint.CheckpointHandler) (Engine, error) {
	return &Solo{local, db, cp}, nil
}
