package server

import "github.com/dbadoy/grinder/server/cft"

type journalObject interface {
	revert(engine cft.Engine) error
}

type insertContract struct {
	key []byte
}

func (i *insertContract) revert(engine cft.Engine) error {
	return engine.Delete(i.key)
}
