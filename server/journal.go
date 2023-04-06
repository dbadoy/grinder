package server

import "github.com/dbadoy/grinder/server/cft"

type journalObject interface {
	Revert(engine cft.Engine) error
}

type insertContract struct {
	key []byte
}

func (i *insertContract) Revert(engine cft.Engine) error {
	return engine.Delete(i.key)
}
