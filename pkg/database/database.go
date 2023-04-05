package database

import "errors"

var (
	ErrAlreadyExist = errors.New("already exist key")
)

type Database interface {
	HealthCheck() error
	Insert(key []byte, data Data) error
	Put(key []byte, data Data) error
	Delete(key []byte) error
	Exist(index string, key []byte) (bool, error)
}

type Data interface {
	// Index must be lower case.
	Index() string
}
