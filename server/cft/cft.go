package cft

import (
	"errors"
	"net"

	"github.com/dbadoy/grinder/pkg/checkpoint"
	"github.com/dbadoy/grinder/pkg/database"
)

var (
	_ Engine = (*CFT)(nil)
)

type CFT struct {
	localnode net.Addr
	srv       server
	db        database.Database
	cp        checkpoint.CheckpointHandler
}

func (c *CFT) Insert(key []byte, data database.Data) error {
	if !c.srv.HasLeaderPermissions() {
		return errors.New("foo")
	}

	return c.db.Insert(key, data)
}

func (c *CFT) Put(key []byte, data database.Data) error {
	if !c.srv.HasLeaderPermissions() {
		return errors.New("foo")
	}

	return c.db.Put(key, data)
}

func (c *CFT) Delete(key []byte) error {
	if !c.srv.HasLeaderPermissions() {
		return errors.New("foo")
	}

	return c.db.Delete(key)
}

func (c *CFT) Exist(index string, key []byte) (bool, error) {
	return c.db.Exist(index, key)
}

func (c *CFT) Checkpoint() uint64 {
	return c.cp.Checkpoint()
}

func (c *CFT) SetCheckpoint(cp uint64) error {
	return c.srv.CommitCheckpoint(cp)
}

func (c *CFT) Increase() error {
	return c.srv.CommitCheckpoint(c.cp.Checkpoint() + 1)
}

func (c *CFT) Decrease() error {
	return c.srv.CommitCheckpoint(c.cp.Checkpoint() - 1)
}

func (c *CFT) LocalNode() net.Addr     { return c.srv.LocalNode() }
func (c *CFT) ClusterLeader() net.Addr { return c.srv.ClusterLeader() }
