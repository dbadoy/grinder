package memdb

import (
	"errors"
	"sync"

	"github.com/dbadoy/grinder/pkg/database"
)

var _ = database.Database(&MemoryDB{})

type MemoryDB struct {
	mu sync.Mutex
	v  map[interface{}]database.Data
}

func New() *MemoryDB {
	return &MemoryDB{v: make(map[interface{}]database.Data)}
}

func (m *MemoryDB) HealthCheck() error {
	return nil
}

func (m *MemoryDB) Insert(key []byte, data database.Data) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.v[key]; ok {
		return database.ErrAlreadyExist
	}

	m.v[key] = data
	return nil
}

func (m *MemoryDB) Put(key []byte, data database.Data) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.v[key] = data
	return nil
}

func (m *MemoryDB) Delete(key []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.v[key]; !ok {
		return errors.New("not exist key")
	}

	delete(m.v, key)
	return nil
}

func (m *MemoryDB) Exist(_ string, key []byte) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.v[key]
	return ok, nil
}
