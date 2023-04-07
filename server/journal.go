package server

import "github.com/dbadoy/grinder/server/cft"

// journalObject has contrasting methods for specific
// behaviors. Stored data related to blockchain rarely
// undergoes modifications, only additions. Therefore,
// journalObject simply implements revert, and does
// not define behaviors such as 'rollback(to previous
// value)'.
type journalObject interface {
	revert(engine cft.Engine) error
}

type insertContract struct {
	key []byte
}

func (i *insertContract) revert(engine cft.Engine) error {
	return engine.Delete(i.key)
}
